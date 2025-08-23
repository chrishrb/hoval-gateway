package pubsub

import "context"

type MessageHandler interface {
	Handle(ctx context.Context, receiverMask uint32, message *Message)
}

type MessageHandlerFunc func(ctx context.Context, receiverMask uint32, message *Message)

func (h MessageHandlerFunc) Handle(ctx context.Context, receiverMask uint32, message *Message) {
	h(ctx, receiverMask, message)
}

type Listener interface {
	Connect(ctx context.Context, handler MessageHandler) (Connection, error)
}

type Connection interface {
	Disconnect(ctx context.Context) error
}
