package hoval_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/datatype"
	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ventilationModeSelection = hoval.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 0,
		DatapointID:    40650,
		DatapointType:  datatype.List,
	}
	example = hoval.Datapoint{
		FunctionGroup:  10,
		FunctionNumber: 0,
		DatapointID:    40650,
		DatapointType:  datatype.S32,
	}
)

func setupSvc() *hoval.HovalService {
	store := hoval.NewDatapointStore()
	store.Add("ventilation-selection", ventilationModeSelection)
	store.Add("example", example)
	return hoval.NewHovalService(store, nil)
}

func TestToTransportMessage(t *testing.T) {
	svc := setupSvc()

	msg := hoval.NewMessage(
		1153,
		1,
		0x40,
		ventilationModeSelection,
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
	svc := setupSvc()

	msg := hoval.NewMessage(
		1153,
		1,
		0x40,
		ventilationModeSelection,
		uint32(12),
	)

	_, err := svc.ToTransportMessage(msg)
	require.Error(t, err)
}

func TestFromTransportMessage(t *testing.T) {
	svc := setupSvc()

	tMsg := transport.Message{
		ID:     0x1fe40801,
		Length: 8,
		Data:   transport.Data{0x1, 0x40, 0x32, 0x0, 0x9e, 0xca, 0x1, 0x0},
	}
	msg, err := svc.FromTransportMessage(tMsg)

	require.NoError(t, err)
	require.NotNil(t, msg)

	assert.Equal(t, uint32(1153), msg.SenderID)
	assert.Equal(t, uint32(1), msg.ReceiverMask)
	assert.Equal(t, hoval.Operation(0x40), msg.OperationID)
	assert.Equal(t, ventilationModeSelection, msg.Datapoint)
	assert.Equal(t, uint8(1), msg.Data)
}

func TestFromTransportMessageInvalidData(t *testing.T) {
	svc := setupSvc()

	// Too long data
	tMsg := transport.Message{
		ID:     0x1fe40801,
		Length: 8,
		Data:   transport.Data{0x3, 0x40, 0x32, 0x0, 0x9e, 0xca, 0x1, 0x0},
	}

	_, err := svc.FromTransportMessage(tMsg)
	require.Error(t, err)
}
