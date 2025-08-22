package hoval

// Operation represents operations in CAN-Message
type Operation uint8

const (
	OperationResponse   Operation = 66
	OperationGetRequest Operation = 64
	OperationSetRequest Operation = 70
)
