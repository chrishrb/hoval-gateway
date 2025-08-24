package can

import (
	"context"
	"errors"

	"github.com/chrishrb/hoval-gateway/transport"
	"go.einride.tech/can/pkg/socketcan"
)

type Consumer struct {
	connectionDetails
}

func NewConsumer(opts ...Opt[Consumer]) *Consumer {
	c := &Consumer{}
	for _, opt := range opts {
		opt(c)
	}
	ensureConsumerDefaults(c)
	return c
}

func ensureConsumerDefaults(c *Consumer) {
	if c.device == "" {
		c.device = "can0"
	}
}

func (c *Consumer) Consume(ctx context.Context, handler transport.MessageHandler) (transport.Connection, error) {
	conn, err := socketcan.DialContext(ctx, "can", c.device)
	if err != nil {
		return nil, err
	}

	recv := socketcan.NewReceiver(conn)
	go func() {
		for recv.Receive() {
			frame := recv.Frame()

			msg := &transport.Message{
				ID:     frame.ID,
				Length: frame.Length,
				Data:   transport.Data(frame.Data),
			}

			handler.Handle(ctx, msg)
		}
	}()

	select {
	case <-ctx.Done():
		return nil, errors.New("timeout waiting for can setup")
	default:
		return conn, nil
	}
}
