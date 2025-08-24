// Can be tested only on linux platform because it requires socketcan
//go:build linux

package can_test

import (
	"context"
	"testing"
	"time"

	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/chrishrb/hoval-gateway/transport/can"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	go_can "go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

func TestNewConsumer(t *testing.T) {
	t.Run("default device", func(t *testing.T) {
		consumer := can.NewConsumer()
		assert.NotNil(t, consumer)
	})

	t.Run("with custom device", func(t *testing.T) {
		consumer := can.NewConsumer(can.WithCANDevice[can.Consumer]("vcan1"))
		assert.NotNil(t, consumer)
	})
}

func TestConsumer_Consume(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create consumer with vcan0 device
	consumer := can.NewConsumer(can.WithCANDevice[can.Consumer]("vcan0"))

	// Channel to collect received messages
	receivedMessages := make(chan *transport.Message, 10)
	var handlerCalled bool

	// Create message handler
	handler := transport.MessageHandlerFunc(func(ctx context.Context, message *transport.Message) {
		handlerCalled = true
		select {
		case receivedMessages <- message:
		case <-ctx.Done():
		}
	})

	// Start consuming
	conn, err := consumer.Consume(ctx, handler)
	require.NoError(t, err)
	require.NotNil(t, conn)
	//nolint:errcheck
	defer conn.Close()

	// Give consumer time to set up
	time.Sleep(100 * time.Millisecond)

	// Create sender connection to send test message
	sendConn, err := socketcan.DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	//nolint:errcheck
	defer sendConn.Close()

	// Send test message
	testFrame := go_can.Frame{
		ID:     0x123,
		Length: 8,
		Data:   go_can.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	tx := socketcan.NewTransmitter(sendConn)
	err = tx.TransmitFrame(ctx, testFrame)
	require.NoError(t, err)

	// Wait for message to be received
	select {
	case receivedMsg := <-receivedMessages:
		assert.True(t, handlerCalled)
		assert.Equal(t, uint32(0x123), receivedMsg.ID)
		assert.Equal(t, uint8(8), receivedMsg.Length)
		expectedData := transport.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
		assert.Equal(t, expectedData, receivedMsg.Data)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

func TestConsumer_ConsumeMultipleMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer := can.NewConsumer(can.WithCANDevice[can.Consumer]("vcan0"))

	// Channel to collect received messages
	receivedMessages := make(chan *transport.Message, 10)

	handler := transport.MessageHandlerFunc(func(ctx context.Context, message *transport.Message) {
		select {
		case receivedMessages <- message:
		case <-ctx.Done():
		}
	})

	conn, err := consumer.Consume(ctx, handler)
	require.NoError(t, err)
	require.NotNil(t, conn)
	//nolint:errcheck
	defer conn.Close()

	// Give consumer time to set up
	time.Sleep(100 * time.Millisecond)

	// Create sender
	sendConn, err := socketcan.DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	//nolint:errcheck
	defer sendConn.Close()

	tx := socketcan.NewTransmitter(sendConn)

	// Send multiple test messages
	testMessages := []go_can.Frame{
		{ID: 0x100, Length: 4, Data: go_can.Data{0x01, 0x02, 0x03, 0x04}},
		{ID: 0x200, Length: 2, Data: go_can.Data{0xAA, 0xBB}},
		{ID: 0x300, Length: 8, Data: go_can.Data{0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88}},
	}

	for _, frame := range testMessages {
		err = tx.TransmitFrame(ctx, frame)
		require.NoError(t, err)
	}

	// Collect received messages
	var received []*transport.Message
	timeout := time.After(2 * time.Second)

	for i := range testMessages {
		select {
		case msg := <-receivedMessages:
			received = append(received, msg)
		case <-timeout:
			t.Fatalf("Timeout waiting for message %d", i+1)
		}
	}

	// Verify all messages were received correctly
	assert.Len(t, received, len(testMessages))

	// Check each message (order might not be preserved)
	receivedIDs := make(map[uint32]*transport.Message)
	for _, msg := range received {
		receivedIDs[msg.ID] = msg
	}

	for _, expectedFrame := range testMessages {
		receivedMsg, exists := receivedIDs[expectedFrame.ID]
		assert.True(t, exists, "Message with ID %x not received", expectedFrame.ID)
		if exists {
			assert.Equal(t, expectedFrame.Length, receivedMsg.Length)
			expectedData := transport.Data{}
			copy(expectedData[:], expectedFrame.Data[:])
			assert.Equal(t, expectedData, receivedMsg.Data)
		}
	}
}

func TestConsumer_Consume_InvalidDevice(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Use non-existent device
	consumer := can.NewConsumer(can.WithCANDevice[can.Consumer]("nonexistent"))

	handler := transport.MessageHandlerFunc(func(ctx context.Context, message *transport.Message) {
		// Handler implementation
	})

	conn, err := consumer.Consume(ctx, handler)
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestConsumer_MessageHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	consumer := can.NewConsumer(can.WithCANDevice[can.Consumer]("vcan0"))

	// Track handler calls
	var handlerCalls []transport.Message
	handlerMutex := make(chan struct{}, 1)
	handlerMutex <- struct{}{} // Initialize mutex

	handler := transport.MessageHandlerFunc(func(ctx context.Context, message *transport.Message) {
		<-handlerMutex // Acquire mutex
		handlerCalls = append(handlerCalls, *message)
		handlerMutex <- struct{}{} // Release mutex
	})

	conn, err := consumer.Consume(ctx, handler)
	require.NoError(t, err)
	require.NotNil(t, conn)
	//nolint:errcheck
	defer conn.Close()

	// Give consumer time to set up
	time.Sleep(100 * time.Millisecond)

	// Create sender
	sendConn, err := socketcan.DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	//nolint:errcheck
	defer sendConn.Close()

	tx := socketcan.NewTransmitter(sendConn)

	// Send a test message
	testFrame := go_can.Frame{
		ID:     0x456,
		Length: 3,
		Data:   go_can.Data{0xDE, 0xAD, 0xBE, 0xEF},
	}

	err = tx.TransmitFrame(ctx, testFrame)
	require.NoError(t, err)

	// Wait for handler to be called
	time.Sleep(500 * time.Millisecond)

	// Verify handler was called with correct message
	<-handlerMutex // Acquire mutex for reading
	assert.Len(t, handlerCalls, 1)
	if len(handlerCalls) > 0 {
		call := handlerCalls[0]
		assert.Equal(t, uint32(0x456), call.ID)
		assert.Equal(t, uint8(3), call.Length)
		expectedData := transport.Data{0xDE, 0xAD, 0xBE, 0xEF}
		assert.Equal(t, expectedData, call.Data)
	}
	handlerMutex <- struct{}{} // Release mutex
}
