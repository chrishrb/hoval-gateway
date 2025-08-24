//go:build linux

package can_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	go_can "go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

func TestNewCAN(t *testing.T) {
	recvConn, err := socketcan.DialContext(t.Context(), "can", "vcan0")
	require.NoError(t, err)

	// Receive message
	recv := socketcan.NewReceiver(recvConn)

	// Channel to receive the frame
	frameChan := make(chan go_can.Frame, 1)

	// Start receiver in goroutine
	go func() {
		for recv.Receive() {
			frameChan <- recv.Frame()
		}
	}()

	frame := go_can.Frame{
		ID:     0x123,
		Length: 8,
		Data:   go_can.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	sendConn, err := socketcan.DialContext(t.Context(), "can", "vcan0")
	require.NoError(t, err)

	// Send message
	tx := socketcan.NewTransmitter(sendConn)
	err = tx.TransmitFrame(t.Context(), frame)
	require.NoError(t, err)

	// Wait for frame or timeout
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	select {
	case receivedFrame := <-frameChan:
		assert.Equal(t, receivedFrame.ID, uint32(0x123))
		assert.Equal(t, receivedFrame.Length, uint8(8))
		assert.Equal(t, receivedFrame.Data, go_can.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})
	case <-ctx.Done():
		t.Fatal("Test timed out waiting for CAN message")
	}
}
