package transport

import "context"

type MessageHandler interface {
	Handle(ctx context.Context, message *Message)
}

type MessageHandlerFunc func(ctx context.Context, message *Message)

func (h MessageHandlerFunc) Handle(ctx context.Context, message *Message) {
	h(ctx, message)
}

type Consumer interface {
	// Connect establishes a connection and subscribes to receive messages
	//
	// Returns either a Connection on success or an error.
	Consume(ctx context.Context, handler MessageHandler) (Connection, error)
}

type Connection interface {
	// Close drops the previously established connection.
	Close() error
}
