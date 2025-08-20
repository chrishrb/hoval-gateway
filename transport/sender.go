package transport

import "context"

// Sender defines the contract for sending messages.
type Sender interface {
	Send(ctx context.Context, message *Message) error
}
