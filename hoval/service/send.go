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

type SendService struct {
	senderID        uint32
	store           store.Engine
	dpProvider      datapoint.DatapointProvider
	transportSender transport.Sender
}

func NewSendService(
	senderID uint32,
	store store.Engine,
	dpProvider datapoint.DatapointProvider,
	transportSender transport.Sender,
) *SendService {
	return &SendService{
		senderID:        senderID,
		store:           store,
		dpProvider:      dpProvider,
		transportSender: transportSender,
	}
}

func (s *SendService) Handle(ctx context.Context, receiverMask uint32, msg *pubsub.Message) {
	// We assume that we only want to SET requests via pubsub API
	operation := hoval.OperationSetRequest

	err := s.Send(ctx, receiverMask, operation, msg)
	if err != nil {
		slog.Error("error handling pubsub send", "error", err)
	}
}

func (s *SendService) Send(ctx context.Context, receiverMask uint32, operationID hoval.Operation, msg *pubsub.Message) error {
	datapoint := s.dpProvider.GetByFunction(msg.FunctionGroup, msg.FunctionNumber, msg.DatapointID)
	if datapoint == nil {
		return fmt.Errorf("datapoint not found: %d/%d/%d", msg.FunctionGroup, msg.FunctionNumber, msg.DatapointID)
	}

	if operationID == hoval.OperationSetRequest && !datapoint.Writable {
		return fmt.Errorf("datapoint is not writable: %d/%d/%d", msg.FunctionGroup, msg.FunctionNumber, msg.DatapointID)
	}

	hovalMsg := &hoval.Message{
		SenderID:     s.senderID,
		ReceiverMask: receiverMask,
		OperationID:  operationID,
		Datapoint:    datapoint,
		Data:         msg.Data,
	}

	transportMsg, err := s.ToTransportMessage(hovalMsg)
	if err != nil {
		return fmt.Errorf("failed to convert message to transport format: %v", err)
	}

	err = s.transportSender.Send(ctx, transportMsg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (s *SendService) ToTransportMessage(msg *hoval.Message) (*transport.Message, error) {
	// Get datapoint from provider
	if msg.Datapoint == nil {
		return nil, fmt.Errorf("datapoint name is required")
	}

	dp := msg.Datapoint

	data := new([8]byte)

	// Add Operation ID
	data[1] = byte(msg.OperationID)

	// Add Datapoint information
	data[2] = byte(dp.FunctionGroup)
	data[3] = byte(dp.FunctionNumber)
	data[4], data[5] = byte(dp.DatapointID>>8), byte(dp.DatapointID)

	// Add data
	d, err := datatype.ToBytes(datatype.Type(dp.TypeName), msg.Data, int(dp.Decimal))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// TODO: data max 2 bytes (16 bit)
	if len(d) >= 1 {
		data[6] = d[0]
	}
	if len(d) == 2 {
		data[7] = d[1]
	}
	if len(d) > 2 {
		return nil, fmt.Errorf("data exceeds maximum length of 2 bytes: %w", hoval.ErrInvalidMessageLength)
	}

	// TODO: handle messages with chunks
	// Special datapoints have more CAN messages because they are sent in chunks.
	// U32, S32, S64 need more than 8 bytes, so the first byte indicates how many messages are needed.
	data[0] = 1

	// Set the frameLen of the data frame
	frameLen := uint8(6 + len(d))

	return &transport.Message{
		ID:     (0x7F << 22) | (msg.SenderID << 11) | msg.ReceiverMask,
		Length: frameLen,
		Data:   *data,
	}, nil
}
