package hoval

import (
	"fmt"

	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/chrishrb/hoval-gateway/transport"
)

type HovalService struct {
	store     *DatapointStore
	canSender *transport.Sender
	// mqttSender  *api.Emiter
}

func NewHovalService(store *DatapointStore, sender *transport.Sender) *HovalService {
	return &HovalService{
		store:     store,
		canSender: sender,
	}
}

func (hs *HovalService) ToTransportMessage(msg *Message) (*transport.Message, error) {
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
		return nil, fmt.Errorf("data exceeds maximum length of 2 bytes: %w", ErrInvalidMessageLength)
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

func (hs *HovalService) FromTransportMessage(msg transport.Message) (*Message, error) {
	// TODO: handle data that is more than 2 bytes
	length := msg.Data[0]
	if length > 2 {
		return nil, fmt.Errorf("data length exceeds 2 bytes: %w", ErrInvalidMessageLength)
	}

	operationID := Operation(msg.Data[1])
	fnGroup := msg.Data[2]
	fnNumber := msg.Data[3]
	datapointID := uint16(msg.Data[4])<<8 | uint16(msg.Data[5])

	datapoint := hs.store.LookupByIdentifier(FunctionGroup(fnGroup), fnNumber, datapointID)
	if datapoint == nil {
		return nil, fmt.Errorf("unknown datapoint: function group %d, function number %d, datapoint ID %d", fnGroup, fnNumber, datapointID)
	}

	data, err := datatype.FromBytes(datapoint.DatapointType, msg.Data[6:6+length], datapoint.DecimalPlaces)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &Message{
		SenderID:     (msg.ID >> 11) & 0x7FF,
		ReceiverMask: msg.ID & 0x7FF,
		OperationID:  operationID,
		Datapoint:    *datapoint,
		Data:         data,
	}, nil
}
