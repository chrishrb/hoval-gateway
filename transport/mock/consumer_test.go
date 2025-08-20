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

func TestConsumer_ReadFromFile(t *testing.T) {
	ctx := context.Background()
	bus := mock.NewMockBus()
	consumer := mock.NewConsumer(bus)

	// Channel to collect received messages
	receivedMessages := make(chan transport.Message, 10)
	done := make(chan struct{})

	// Start a goroutine to collect messages from the bus
	go func() {
		defer close(done)
		for {
			select {
			case msg := <-bus.Bus:
				receivedMessages <- msg
			case <-ctx.Done():
				return
			case <-time.After(100 * time.Millisecond):
				// No more messages expected
				return
			}
		}
	}()

	// Read from test file
	err := consumer.ReadFromFile(ctx, "testfiles/hoval_data_1.log")
	require.NoError(t, err)

	// Wait for processing to complete
	<-done
	close(receivedMessages)

	// Collect all received messages
	var messages []transport.Message
	for msg := range receivedMessages {
		t.Log(msg)
		messages = append(messages, msg)
	}

	// Verify we received the expected number of messages
	assert.GreaterOrEqual(t, len(messages), 1, "should receive at least one message from test file")

	// Verify message format (assuming test file contains valid candump format)
	if len(messages) > 0 {
		msg := messages[0]
		assert.Greater(t, msg.ID, uint32(0), "message ID should be greater than 0")
		assert.LessOrEqual(t, msg.Length, uint8(8), "message length should not exceed 8 bytes")
		// assert.Equal(t, int(msg.Length), len(msg.Data), "message length should match data length")
	}
}

func TestConsumer_ReadFromFile_NonExistentFile(t *testing.T) {
	ctx := context.Background()
	bus := mock.NewMockBus()
	consumer := mock.NewConsumer(bus)

	err := consumer.ReadFromFile(ctx, "nonexistent/file.log")
	assert.Error(t, err, "should return error for non-existent file")
}

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
