package service

import (
	"context"
	"fmt"

	"github.com/chrishrb/hoval-gateway/config"
	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/chrishrb/hoval-gateway/store"
	"github.com/chrishrb/hoval-gateway/transport"
)

type SendService struct {
	store      store.Engine
	dpProvider config.DatapointProvider
	canSender  transport.Sender
}

func NewSendService(store store.Engine, dpProvider config.DatapointProvider, sender transport.Sender) *SendService {
	return &SendService{
		store:      store,
		dpProvider: dpProvider,
		canSender:  sender,
	}
}

func (s *SendService) Handle(ctx context.Context, message *transport.Message) {
}

func (s *SendService) ToTransportMessage(msg *hoval.Message) (*transport.Message, error) {
	// Get datapoint from provider
	if msg.DatapointName == nil {
		return nil, fmt.Errorf("datapoint name is required")
	}

	datapoint := s.dpProvider.GetByName(*msg.DatapointName)
	if datapoint == nil {
		return nil, fmt.Errorf("datapoint not found: %s", *msg.DatapointName)
	}

	data := new([8]byte)

	// Add Operation ID
	data[1] = byte(msg.OperationID)

	// Add Datapoint information
	data[2] = byte(datapoint.FunctionGroup)
	data[3] = byte(datapoint.FunctionNumber)
	data[4], data[5] = byte(datapoint.DatapointID>>8), byte(datapoint.DatapointID)

	// Add data
	d, err := datatype.ToBytes(datatype.Type(datapoint.TypeName), msg.Data, int(datapoint.Decimal))
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
