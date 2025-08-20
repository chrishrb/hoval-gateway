package hoval

import (
	"errors"
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
	Datapoint    Datapoint
	Data         any
}

func NewMessage(
	senderID,
	receiverMask uint32,
	operationID Operation,
	datapoint Datapoint,
	data any,
) *Message {
	return &Message{
		SenderID:     senderID,
		ReceiverMask: receiverMask,
		OperationID:  operationID,
		Datapoint:    datapoint,
		Data:         data,
	}
}
