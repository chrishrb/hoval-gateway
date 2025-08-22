package service_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/service"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProduceSvc() *service.SendService {
	dpProvider := &DatapointProviderMock{}
	store := inmemory.NewStore(nil)
	return service.NewSendService(store, dpProvider, nil)
}

func TestToTransportMessage(t *testing.T) {
	svc := setupProduceSvc()

	dp := &hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 0,
		DatapointID:    40650,
	}

	msg := hoval.NewMessage(
		1153,
		1,
		0x40,
		dp,
		uint8(1),
	)

	tMsg, err := svc.ToTransportMessage(msg)
	require.NoError(t, err)
	require.NotNil(t, tMsg)

	assert.Equal(t, uint32(0x1fe40801), tMsg.ID)
	assert.Equal(t, uint8(7), tMsg.Length)
	assert.Equal(t, transport.Data{0x1, 0x40, 0x32, 0x0, 0x9e, 0xca, 0x1, 0x0}, tMsg.Data)
}

func TestToTransportMessageInvalidData(t *testing.T) {
	svc := setupProduceSvc()

	dp := &hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 0,
		DatapointID:    40650,
	}

	msg := hoval.NewMessage(
		1153,
		1,
		0x40,
		dp,
		uint32(12),
	)

	_, err := svc.ToTransportMessage(msg)
	require.Error(t, err)
}
