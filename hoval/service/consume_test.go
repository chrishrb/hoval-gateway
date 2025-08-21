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
		Length: 8,
		Data:   transport.Data{0x1, 0x40, 0x32, 0x0, 0x9e, 0xca, 0x1, 0x0},
	}
	msg, err := svc.FromTransportMessage(tMsg)

	require.NoError(t, err)
	require.NotNil(t, msg)

	assert.Equal(t, uint32(1153), msg.SenderID)
	assert.Equal(t, uint32(1), msg.ReceiverMask)
	assert.Equal(t, hoval.Operation(0x40), msg.OperationID)
	assert.Equal(t, "TestDatapoint", *msg.DatapointName)
	assert.Equal(t, uint8(1), msg.Data)
}

func TestFromTransportMessageInvalidData(t *testing.T) {
	svc := setupConsumeSvc()

	// Too long data
	tMsg := transport.Message{
		ID:     0x1fe40801,
		Length: 8,
		Data:   transport.Data{0x3, 0x40, 0x32, 0x0, 0x9e, 0xca, 0x1, 0x0},
	}

	_, err := svc.FromTransportMessage(tMsg)
	require.Error(t, err)
}
