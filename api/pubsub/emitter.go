package pubsub

import (
	"context"
)

type Emitter interface {
	Emit(ctx context.Context, receiverMask uint32, message *Message) error
}

// EmitterFunc allows a plain function to be used as an Emitter
type EmitterFunc func(ctx context.Context, receiverMask uint32, message *Message) error

func (e EmitterFunc) Emit(ctx context.Context, receiverMask uint32, message *Message) error {
	return e(ctx, receiverMask, message)
}
