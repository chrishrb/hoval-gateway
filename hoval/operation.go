package hoval

// Operation represents operations in CAN-Message
type Operation uint8

const (
	OperationResponse   Operation = 0x42
	OperationGetRequest Operation = 0x40
	OperationSetRequest Operation = 0x46
)
