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

func TestNewSender(t *testing.T) {
	t.Run("default device", func(t *testing.T) {
		sender := can.NewSender()
		assert.NotNil(t, sender)
	})

	t.Run("with custom device", func(t *testing.T) {
		sender := can.NewSender(can.WithCANDevice[can.Sender]("vcan1"))
		assert.NotNil(t, sender)
	})
}

func TestSender_Send(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create sender with vcan0 device
	sender := can.NewSender(can.WithCANDevice[can.Sender]("vcan0"))

	// Create receiver to verify the message was sent
	receiverConn, err := socketcan.DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	//nolint:errcheck
	defer receiverConn.Close()

	receiver := socketcan.NewReceiver(receiverConn)

	// Create test message
	testMessage := &transport.Message{
		ID:     0x123,
		Length: 8,
		Data:   transport.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	// Channel to receive frames
	receivedFrames := make(chan go_can.Frame, 1)

	// Start receiving in goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if receiver.Receive() {
					frame := receiver.Frame()
					select {
					case receivedFrames <- frame:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Give receiver time to set up
	time.Sleep(100 * time.Millisecond)

	// Send the message
	err = sender.Send(ctx, testMessage)
	require.NoError(t, err)

	// Verify the frame was received
	select {
	case frame := <-receivedFrames:
		assert.Equal(t, uint32(0x123), frame.ID)
		assert.Equal(t, uint8(8), frame.Length)
		assert.Equal(t, go_can.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, frame.Data)
		assert.False(t, frame.IsRemote)
		assert.True(t, frame.IsExtended)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for sent message")
	}
}

func TestSender_SendMultipleMessages(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sender := can.NewSender(can.WithCANDevice[can.Sender]("vcan0"))

	// Create receiver
	receiverConn, err := socketcan.DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	//nolint:errcheck
	defer receiverConn.Close()

	receiver := socketcan.NewReceiver(receiverConn)

	// Test messages
	testMessages := []*transport.Message{
		{ID: 0x100, Length: 4, Data: transport.Data{0x01, 0x02, 0x03, 0x04}},
		{ID: 0x200, Length: 2, Data: transport.Data{0xAA, 0xBB}},
		{ID: 0x300, Length: 8, Data: transport.Data{0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88}},
	}

	// Channel to receive frames
	receivedFrames := make(chan go_can.Frame, len(testMessages))

	// Start receiving
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if receiver.Receive() {
					frame := receiver.Frame()
					select {
					case receivedFrames <- frame:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Give receiver time to set up
	time.Sleep(100 * time.Millisecond)

	// Send all messages
	for _, msg := range testMessages {
		err = sender.Send(ctx, msg)
		require.NoError(t, err)
	}

	// Collect received frames
	var received []go_can.Frame
	timeout := time.After(2 * time.Second)

	for i := range testMessages {
		select {
		case frame := <-receivedFrames:
			received = append(received, frame)
		case <-timeout:
			t.Fatalf("Timeout waiting for message %d", i+1)
		}
	}

	// Verify all messages were received correctly
	assert.Len(t, received, len(testMessages))

	// Check each message (order might not be preserved)
	receivedIDs := make(map[uint32]go_can.Frame)
	for _, frame := range received {
		receivedIDs[frame.ID] = frame
	}

	for _, expectedMsg := range testMessages {
		receivedFrame, exists := receivedIDs[expectedMsg.ID]
		assert.True(t, exists, "Message with ID %x not received", expectedMsg.ID)
		if exists {
			assert.Equal(t, expectedMsg.Length, receivedFrame.Length)
			expectedData := go_can.Data{}
			copy(expectedData[:], expectedMsg.Data[:])
			assert.Equal(t, expectedData, receivedFrame.Data)
			assert.False(t, receivedFrame.IsRemote)
			assert.True(t, receivedFrame.IsExtended)
		}
	}
}

func TestSender_Send_InvalidDevice(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Use non-existent device
	sender := can.NewSender(can.WithCANDevice[can.Sender]("nonexistent"))

	testMessage := &transport.Message{
		ID:     0x123,
		Length: 4,
		Data:   transport.Data{0x01, 0x02, 0x03, 0x04},
	}

	err := sender.Send(ctx, testMessage)
	assert.Error(t, err)
}

func TestSender_Send_ContextCancellation(t *testing.T) {
	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sender := can.NewSender(can.WithCANDevice[can.Sender]("vcan0"))

	testMessage := &transport.Message{
		ID:     0x123,
		Length: 4,
		Data:   transport.Data{0x01, 0x02, 0x03, 0x04},
	}

	err := sender.Send(ctx, testMessage)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestSender_ConnectionReuse(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sender := can.NewSender(can.WithCANDevice[can.Sender]("vcan0"))

	// Create receiver
	receiverConn, err := socketcan.DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	//nolint:errcheck
	defer receiverConn.Close()

	receiver := socketcan.NewReceiver(receiverConn)

	// Channel to receive frames
	receivedFrames := make(chan go_can.Frame, 2)

	// Start receiving
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if receiver.Receive() {
					frame := receiver.Frame()
					select {
					case receivedFrames <- frame:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Give receiver time to set up
	time.Sleep(100 * time.Millisecond)

	// Send first message
	msg1 := &transport.Message{
		ID:     0x100,
		Length: 4,
		Data:   transport.Data{0x01, 0x02, 0x03, 0x04},
	}
	err = sender.Send(ctx, msg1)
	require.NoError(t, err)

	// Send second message (should reuse connection)
	msg2 := &transport.Message{
		ID:     0x200,
		Length: 4,
		Data:   transport.Data{0x05, 0x06, 0x07, 0x08},
	}
	err = sender.Send(ctx, msg2)
	require.NoError(t, err)

	// Verify both messages were received
	var receivedCount int
	timeout := time.After(2 * time.Second)

	for receivedCount < 2 {
		select {
		case frame := <-receivedFrames:
			receivedCount++
			switch frame.ID {
			case 0x100:
				assert.Equal(t, go_can.Data{0x01, 0x02, 0x03, 0x04}, frame.Data)
			case 0x200:
				assert.Equal(t, go_can.Data{0x05, 0x06, 0x07, 0x08}, frame.Data)
			default:
				t.Fatalf("Unexpected frame ID: %x", frame.ID)
			}
		case <-timeout:
			t.Fatalf("Timeout waiting for messages, received %d of 2", receivedCount)
		}
	}

	assert.Equal(t, 2, receivedCount)
}

func TestSender_ConcurrentSend(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sender := can.NewSender(can.WithCANDevice[can.Sender]("vcan0"))

	// Create receiver
	receiverConn, err := socketcan.DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	//nolint:errcheck
	defer receiverConn.Close()

	receiver := socketcan.NewReceiver(receiverConn)

	// Number of concurrent senders
	numSenders := 10
	receivedFrames := make(chan go_can.Frame, numSenders)

	// Start receiving
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if receiver.Receive() {
					frame := receiver.Frame()
					select {
					case receivedFrames <- frame:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Give receiver time to set up
	time.Sleep(100 * time.Millisecond)

	// Send messages concurrently
	errChan := make(chan error, numSenders)
	for i := range numSenders {
		go func(id int) {
			msg := &transport.Message{
				ID:     uint32(0x100 + id),
				Length: 4,
				Data:   transport.Data{byte(id), byte(id + 1), byte(id + 2), byte(id + 3)},
			}
			errChan <- sender.Send(ctx, msg)
		}(i)
	}

	// Check all sends completed without error
	for range numSenders {
		err := <-errChan
		assert.NoError(t, err)
	}

	// Verify all messages were received
	receivedIDs := make(map[uint32]bool)
	timeout := time.After(3 * time.Second)

	for len(receivedIDs) < numSenders {
		select {
		case frame := <-receivedFrames:
			receivedIDs[frame.ID] = true
		case <-timeout:
			t.Fatalf("Timeout waiting for messages, received %d of %d", len(receivedIDs), numSenders)
		}
	}

	// Verify all expected IDs were received
	for i := range numSenders {
		expectedID := uint32(0x100 + i)
		assert.True(t, receivedIDs[expectedID], "Message with ID %x not received", expectedID)
	}
}

func TestSender_FrameProperties(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sender := can.NewSender(can.WithCANDevice[can.Sender]("vcan0"))

	// Create receiver
	receiverConn, err := socketcan.DialContext(ctx, "can", "vcan0")
	require.NoError(t, err)
	//nolint:errcheck
	defer receiverConn.Close()

	receiver := socketcan.NewReceiver(receiverConn)

	// Channel to receive frames
	receivedFrames := make(chan go_can.Frame, 1)

	// Start receiving
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if receiver.Receive() {
					frame := receiver.Frame()
					select {
					case receivedFrames <- frame:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	// Give receiver time to set up
	time.Sleep(100 * time.Millisecond)

	// Test message with specific properties
	testMessage := &transport.Message{
		ID:     0x1FFFFFFF, // Max 29-bit extended ID
		Length: 0,          // Empty message
		Data:   transport.Data{},
	}

	err = sender.Send(ctx, testMessage)
	require.NoError(t, err)

	// Verify frame properties
	select {
	case frame := <-receivedFrames:
		assert.Equal(t, uint32(0x1FFFFFFF), frame.ID)
		assert.Equal(t, uint8(0), frame.Length)
		assert.False(t, frame.IsRemote, "Frame should not be remote")
		assert.True(t, frame.IsExtended, "Frame should be extended")
		// Data should be zero for empty message
		assert.Equal(t, go_can.Data{}, frame.Data)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for sent message")
	}
}
