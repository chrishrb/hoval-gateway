package hoval

import (
	"errors"

	"github.com/chrishrb/hoval-gateway/hoval/datapoint"
)

var (
	ErrInvalidMessageLength = errors.New("invalid message length")
)

// Message represents a Hoval CAN message structure
type Message struct {
	// SenderID is the ID of the sender device
	SenderID uint32
	// ReceiverMask is a bitmask of receiver IDs that should receive this message
	ReceiverMask uint32
	OperationID  Operation
	Datapoint    *datapoint.Datapoint
	Data         float64
}

func NewMessage(
	senderID,
	receiverMask uint32,
	operationID Operation,
	datapoint *datapoint.Datapoint,
	data float64,
) *Message {
	return &Message{
		SenderID:     senderID,
		ReceiverMask: receiverMask,
		OperationID:  operationID,
		Datapoint:    datapoint,
		Data:         data,
	}
}
