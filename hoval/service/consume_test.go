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

func setupConsumeSvc() *service.ConsumeService {
	dpProvider := &DatapointProviderMock{}
	store := inmemory.NewStore(nil)
	return service.NewConsumeService(store, dpProvider)
}

func TestFromTransportMessage(t *testing.T) {
	svc := setupConsumeSvc()

	tMsg := transport.Message{
		ID:     0x1fe40801,
		Length: 7,
		Data:   transport.Data{0x1, 0x42, 0x32, 0x0, 0x9e, 0xca, 0x1},
	}
	msg, err := svc.FromTransportMessage(tMsg)

	require.NoError(t, err)
	require.NotNil(t, msg)

	assert.Equal(t, uint32(1153), msg.SenderID)
	assert.Equal(t, uint32(1), msg.ReceiverMask)
	assert.Equal(t, hoval.Operation(0x42), msg.OperationID)
	assert.Equal(t, uint8(50), msg.Datapoint.FunctionGroup)
	assert.Equal(t, uint8(0), msg.Datapoint.FunctionNumber)
	assert.Equal(t, uint16(40650), msg.Datapoint.DatapointID)
	assert.Equal(t, uint8(1), msg.Data)
}

func TestFromTransportMessageInvalidData(t *testing.T) {
	svc := setupConsumeSvc()

	// Too long for LIST datatype
	tMsg := transport.Message{
		ID:     0x1fe40801,
		Length: 8,
		Data:   transport.Data{0x1, 0x42, 0x32, 0x0, 0x9e, 0xca, 0x0, 0x1},
	}

	_, err := svc.FromTransportMessage(tMsg)
	require.Error(t, err)
}
