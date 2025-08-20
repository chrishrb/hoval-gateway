package mock

import (
	"context"
	"errors"

	"github.com/chrishrb/hoval-gateway/transport"
)

type Consumer struct {
	conn *MockBus
}

func NewConsumer(bus *MockBus) *Consumer {
	return &Consumer{
		conn: bus,
	}
}

func (c *Consumer) Consume(ctx context.Context, handler transport.MessageHandler) (transport.Connection, error) {
	go func() {
		for {
			msg := <-c.conn.Bus
			handler.Handle(ctx, &msg)
		}
	}()

	select {
	case <-ctx.Done():
		return nil, errors.New("timeout waiting for mock setup")
	default:
		return c.conn, nil
	}
}
