package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chrishrb/hoval-gateway/config"
	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/chrishrb/hoval-gateway/store"
	"github.com/chrishrb/hoval-gateway/transport"
)

type ConsumeService struct {
	store      store.Engine
	dpProvider config.DatapointProvider
	// mqttSender  *api.Emiter
}

func NewConsumeService(store store.Engine, dpProvider config.DatapointProvider) *ConsumeService {
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

	if hovalMsg.DatapointName == nil {
		return
	}

	slog.Debug("received hoval message", "message", hovalMsg)

	// TODO: handle the hoval message, e.g., publish to MQTT or process further
}

func (s *ConsumeService) FromTransportMessage(msg transport.Message) (*hoval.Message, error) {
	// Special datapoints have more CAN messages because they are sent in chunks.
	// U32, S32, S64 need more than 8 bytes, so the first byte indicates how many messages are needed.
	// noOfMsg := msg.Data[0]

	operationID := hoval.Operation(msg.Data[1])
	if operationID != hoval.OperationResponse {
		return nil, nil
	}

	fnGroup := msg.Data[2]
	fnNumber := msg.Data[3]
	datapointID := uint16(msg.Data[4])<<8 | uint16(msg.Data[5])

	datapoint := s.dpProvider.GetByFunction(fnGroup, fnNumber, datapointID)
	if datapoint == nil {
		slog.Error("unknown datapoint", "functionGroup", fnGroup, "functionNumber", fnNumber, "datapointID", datapointID)

		return &hoval.Message{
			SenderID:     (msg.ID >> 11) & 0x7FF,
			ReceiverMask: msg.ID & 0x7FF,
			OperationID:  operationID,
		}, nil
	}

	slog.Debug("datapoint found",
		"operationID", operationID,
		"datapointName", datapoint.DatapointName,
		"datapointType", datapoint.TypeName,
		"data", fmt.Sprintf("0x%x", msg.Data),
	)

	data, err := datatype.FromBytes(datatype.Type(datapoint.TypeName), msg.Data[6:msg.Length], int(datapoint.Decimal))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	dpName := datapoint.DatapointName

	return &hoval.Message{
		SenderID:      (msg.ID >> 11) & 0x7FF,
		ReceiverMask:  msg.ID & 0x7FF,
		OperationID:   operationID,
		DatapointName: &dpName,
		Data:          data,
	}, nil
}
