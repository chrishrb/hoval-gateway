package service

import (
	"context"
	"fmt"

	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/chrishrb/hoval-gateway/transport"
)

type SendService struct {
	store     *hoval.DatapointStore
	canSender transport.Sender
}

func NewSendService(store *hoval.DatapointStore, sender transport.Sender) *SendService {
	return &SendService{
		store:     store,
		canSender: sender,
	}
}

func (s *SendService) Handle(ctx context.Context, message *transport.Message) {
}

func (s *SendService) ToTransportMessage(msg *hoval.Message) (*transport.Message, error) {
	data := new([8]byte)

	// Add Operation ID
	data[1] = byte(msg.OperationID)

	// Add Datapoint information
	dp := msg.Datapoint
	data[2] = byte(dp.FunctionGroup)
	data[3] = byte(dp.FunctionNumber)
	data[4], data[5] = byte(dp.DatapointID>>8), byte(dp.DatapointID)

	// Add data
	d, err := datatype.ToBytes(dp.DatapointType, msg.Data, dp.DecimalPlaces)
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

	// Add message len of complete CAN message
	data[0] = uint8(len(d))

	// Set the frameLen of the data frame
	frameLen := uint8(len(data))

	return &transport.Message{
		ID:     (0x7F << 22) | (msg.SenderID << 11) | msg.ReceiverMask,
		Length: frameLen,
		Data:   *data,
	}, nil
}
