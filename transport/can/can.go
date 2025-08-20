package can

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.einride.tech/can/pkg/candevice"
)

func NewCAN(t *testing.T) *candevice.Device {
	d, _ := candevice.New("vcan0")
	err := d.SetBitrate(250000)
	require.NoError(t, err)

	err = d.SetUp()
	require.NoError(t, err)

	defer d.SetDown()

	return d
}
