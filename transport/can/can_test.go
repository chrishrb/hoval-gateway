//go:build linux

package can_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	go_can "go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

func TestNewCAN(t *testing.T) {
	conn, _ := socketcan.DialContext(t.Context(), "can", "vcan0")

	frame := go_can.Frame{
		ID:     0x123,
		Length: 8,
		Data:   go_can.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	// Send message
	tx := socketcan.NewTransmitter(conn)
	_ = tx.TransmitFrame(t.Context(), frame)

	// Create context with timeout for receiving
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	// Receive message
	recv := socketcan.NewReceiver(conn)

	// Channel to receive the frame
	frameChan := make(chan go_can.Frame, 1)
	
	// Start receiver in goroutine
	go func() {
		for recv.Receive() {
			frameChan <- recv.Frame()
			return // Only receive one frame for this test
		}
	}()

	// Wait for frame or timeout
	select {
	case receivedFrame := <-frameChan:
		assert.Equal(t, receivedFrame.ID, uint32(0x123))
		assert.Equal(t, receivedFrame.Length, uint8(8))
		assert.Equal(t, receivedFrame.Data, go_can.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})
	case <-ctx.Done():
		t.Fatal("Test timed out waiting for CAN message")
	}
}
