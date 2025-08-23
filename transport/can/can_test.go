//go:build linux

package can_test

import (
	"context"
	"testing"
	"time"

	"github.com/chrishrb/hoval-gateway/transport/can"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	go_can "go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

func TestNewCAN(t *testing.T) {
	d := can.NewCAN(t)

	up, err := d.IsUp()
	require.NoError(t, err)
	assert.True(t, up)

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
	for recv.Receive() && ctx.Err() == nil {
		frame := recv.Frame()

		assert.Equal(t, frame.ID, uint32(0x123))
		assert.Equal(t, frame.Length, uint8(8))
		assert.Equal(t, frame.Data, go_can.Data{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08})

		break // Stop after receiving the first message
	}
	
	// Check if we timed out
	if ctx.Err() == context.DeadlineExceeded {
		t.Fatal("Test timed out waiting for CAN message")
	}
}
