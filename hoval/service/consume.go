package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chrishrb/hoval-gateway/api/pubsub"
	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/datapoint"
	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/chrishrb/hoval-gateway/store"
	"github.com/chrishrb/hoval-gateway/transport"
)

type ConsumeService struct {
	store         store.Engine
	dpProvider    datapoint.DatapointProvider
	pubsubEmitter pubsub.Emitter
}

func NewConsumeService(store store.Engine, dpProvider datapoint.DatapointProvider) *ConsumeService {
	return &ConsumeService{
		store:      store,
		dpProvider: dpProvider,
	}
}

func (s *ConsumeService) Handle(ctx context.Context, message *transport.Message) {
	if message == nil {
		slog.Debug("received nil message, ignoring")
		return
	}

	hovalMsg, err := s.FromTransportMessage(*message)
	if err != nil {
		slog.Error("failed to convert transport message to hoval message", "error", err)
		return
	}
	if hovalMsg == nil {
		return
	}

	s.store.SetDevice(hovalMsg.SenderID, &hoval.Device{Address: hovalMsg.SenderID})

	if hovalMsg.Datapoint == nil {
		return
	}

	slog.Debug("received hoval message", "message", hovalMsg)

	err = s.pubsubEmitter.Emit(ctx, hovalMsg.ReceiverMask, &pubsub.Message{
		FunctionGroup:  hovalMsg.Datapoint.FunctionGroup,
		FunctionNumber: hovalMsg.Datapoint.FunctionNumber,
		DatapointID:    hovalMsg.Datapoint.DatapointID,
		Data:           hovalMsg.Data,
	})
	if err != nil {
		slog.Error("failed to emit message to pubsub api", "error", err)
	}
}

func (s *ConsumeService) FromTransportMessage(msg transport.Message) (*hoval.Message, error) {
	if msg.Length < 6 {
		return nil, nil
	}

	// Special datapoints have more CAN messages because they are sent in chunks.
	// U32, S32, S64 need more than 8 bytes, so the first byte indicates how many messages are needed.
	// noOfMsg := msg.Data[0]

	operationID := hoval.Operation(msg.Data[1])
	// if operationID != hoval.OperationResponse {
	// 	slog.Info("skipping messages with operationID not OperationResponse", "operationID", operationID)
	// 	return nil, nil
	// }

	fnGroup := msg.Data[2]
	fnNumber := msg.Data[3]
	datapointID := uint16(msg.Data[4])<<8 | uint16(msg.Data[5])

	dp := s.dpProvider.GetByFunction(fnGroup, fnNumber, datapointID)
	if dp == nil {
		slog.Error("unknown datapoint",
			"operationID", operationID,
			"functionGroup", fnGroup,
			"functionNumber", fnNumber,
			"datapointID", datapointID,
			"data", fmt.Sprintf("0x%x", msg.Data[6:msg.Length]),
		)

		return &hoval.Message{
			SenderID:     (msg.ID >> 11) & 0x7FF,
			ReceiverMask: msg.ID & 0x7FF,
			OperationID:  operationID,
		}, nil
	}

	slog.Info("datapoint found",
		"operationID", operationID,
		"datapointName", dp.DatapointName,
		"datapointType", dp.TypeName,
		"data", fmt.Sprintf("0x%x", msg.Data[6:msg.Length]),
	)

	data, err := datatype.FromBytes(datatype.Type(dp.TypeName), msg.Data[6:msg.Length], int(dp.Decimal))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &hoval.Message{
		SenderID:     (msg.ID >> 11) & 0x7FF,
		ReceiverMask: msg.ID & 0x7FF,
		OperationID:  operationID,
		Datapoint: &datapoint.Datapoint{
			FunctionGroup:  fnGroup,
			FunctionNumber: fnNumber,
			DatapointID:    datapointID,
		},
		Data: data,
	}, nil
}
