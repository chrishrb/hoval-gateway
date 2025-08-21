package service_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/chrishrb/hoval-gateway/hoval/service"
	"github.com/chrishrb/hoval-gateway/store"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProduceSvc() (*service.SendService, store.Engine) {
	ventilationModeSelection := hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 0,
		DatapointID:    40650,
		DatapointType:  datatype.List,
	}
	example := hoval.Datapoint{
		FunctionGroup:  10,
		FunctionNumber: 0,
		DatapointID:    40650,
		DatapointType:  datatype.S32,
	}

	store := inmemory.NewStore(nil)
	store.SetDatapoint("ventilation-selection", &ventilationModeSelection)
	store.SetDatapoint("example", &example)
	return service.NewSendService(store, nil), store
}

func TestToTransportMessage(t *testing.T) {
	svc, store := setupProduceSvc()

	ventilationModeSelection := store.LookupDatapointByName("ventilation-selection")
	require.NotNil(t, ventilationModeSelection)

	msg := hoval.NewMessage(
		1153,
		1,
		0x40,
		*ventilationModeSelection,
		uint8(1),
	)

	tMsg, err := svc.ToTransportMessage(msg)
	require.NoError(t, err)
	require.NotNil(t, tMsg)

	assert.Equal(t, uint32(0x1fe40801), tMsg.ID)
	assert.Equal(t, uint8(8), tMsg.Length)
	assert.Equal(t, transport.Data{0x1, 0x40, 0x32, 0x0, 0x9e, 0xca, 0x1, 0x0}, tMsg.Data)
}

func TestToTransportMessageInvalidData(t *testing.T) {
	svc, store := setupProduceSvc()

	ventilationModeSelection := store.LookupDatapointByName("ventilation-selection")
	require.NotNil(t, ventilationModeSelection)

	msg := hoval.NewMessage(
		1153,
		1,
		0x40,
		*ventilationModeSelection,
		uint32(12),
	)

	_, err := svc.ToTransportMessage(msg)
	require.Error(t, err)
}
