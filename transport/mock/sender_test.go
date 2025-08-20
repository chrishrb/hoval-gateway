package mock_test

import (
	"context"
	"testing"
	"time"

	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/chrishrb/hoval-gateway/transport/mock"
	"github.com/stretchr/testify/assert"
)

func TestSender(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bus := mock.NewMockBus()
	sender := mock.NewSender(bus)
	defer bus.Close()

	// Create a message
	data := transport.Data{1, 2, 3, 4, 5, 6, 7, 8}
	message := transport.NewMessage(2, 8, data)

	receivedMsgCh := make(chan struct{})
	go func() {
		for msg := range bus.Bus {
			assert.Equal(t, message.ID, msg.ID)
			assert.Equal(t, message.Length, msg.Length)
			assert.Equal(t, message.Data, msg.Data)
			receivedMsgCh <- struct{}{}
		}
	}()

	err := sender.Send(ctx, message)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	select {
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test to complete")
	case <-receivedMsgCh:
		// do nothing
	}
}
