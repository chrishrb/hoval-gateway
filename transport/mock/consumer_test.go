package mock_test

import (
	"context"
	"testing"
	"time"

	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/chrishrb/hoval-gateway/transport/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bus := mock.NewMockBus()
	consumer := mock.NewConsumer(bus)

	// Create a message
	data := transport.Data{1, 2, 3, 4, 5, 6, 7, 8}
	message := transport.NewMessage(1, 8, data)

	// Check if message was received
	receivedMsgCh := make(chan struct{})
	handler := func(ctx context.Context, msg *transport.Message) {
		assert.Equal(t, message.ID, msg.ID)
		assert.Equal(t, message.Length, msg.Length)
		assert.Equal(t, message.Data, msg.Data)
		receivedMsgCh <- struct{}{}
	}

	conn, err := consumer.Consume(ctx, transport.MessageHandlerFunc(handler))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() {
		if conn != nil {
			err := conn.Close()
			require.NoError(t, err)
		}
	}()

	// Send message
	bus.Bus <- *message

	// wait for message to be received / timeout
	select {
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test to complete")
	case <-receivedMsgCh:
		// do nothing
	}
}
