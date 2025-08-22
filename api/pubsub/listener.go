package pubsub

import "context"

type MessageHandler interface {
	Handle(ctx context.Context, functionGroup, functionNumber uint8, datapointID uint16, message *Message)
}

type MessageHandlerFunc func(ctx context.Context, functionGroup, functionNumber uint8, datapointID uint16, message *Message)

func (h MessageHandlerFunc) Handle(ctx context.Context, functionGroup, functionNumber uint8, datapointID uint16, message *Message) {
	h(ctx, functionGroup, functionNumber, datapointID, message)
}

type Listener interface {
	Connect(ctx context.Context, handler MessageHandler) (Connection, error)
}

type Connection interface {
	Disconnect(ctx context.Context) error
}
