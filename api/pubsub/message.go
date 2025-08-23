package pubsub

type Message struct {
	FunctionGroup  uint8
	FunctionNumber uint8
	DatapointID    uint16
	Data           float64
}
