package can

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.einride.tech/can/pkg/candevice"
)

func NewCAN(t *testing.T) *candevice.Device {
	d, err := candevice.New("vcan0")
	require.NoError(t, err)
	// err = d.SetBitrate(250000)
	// require.NoError(t, err)

	err = d.SetUp()
	require.NoError(t, err)

	//nolint:errcheck
	defer d.SetDown()

	return d
}
