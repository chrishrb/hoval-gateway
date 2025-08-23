package service_test

import (
	"context"
	"testing"

	"github.com/chrishrb/hoval-gateway/api/pubsub"
	"github.com/chrishrb/hoval-gateway/hoval"
	"github.com/chrishrb/hoval-gateway/hoval/datapoint"
	"github.com/chrishrb/hoval-gateway/hoval/service"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockSender struct {
	Messages []*transport.Message
}

func NewMockSender() *MockSender {
	return &MockSender{
		Messages: []*transport.Message{},
	}
}

func (m *MockSender) Send(ctx context.Context, msg *transport.Message) error {
	m.Messages = append(m.Messages, msg)
	return nil
}

func setupProduceSvc() (*service.SendService, *MockSender) {
	senderID := uint32(1)
	mockSender := NewMockSender()
	dpProvider := &DatapointProviderMock{}
	store := inmemory.NewStore(nil)
	return service.NewSendService(senderID, store, dpProvider, mockSender), mockSender
}

func TestHandle(t *testing.T) {
	svc, sender := setupProduceSvc()

	msg := &pubsub.Message{
		FunctionGroup:  50,
		FunctionNumber: 0,
		DatapointID:    40650,
		Data:           2,
	}

	svc.Handle(context.Background(), 2, msg)

	assert.Len(t, sender.Messages, 1)
	assert.Equal(t, uint32(0x1fc00802), sender.Messages[0].ID)
	assert.Equal(t, uint8(7), sender.Messages[0].Length)
	assert.Equal(t, transport.Data{0x1, 0x46, 0x32, 0x0, 0x9e, 0xca, 0x2, 0x0}, sender.Messages[0].Data)
}

func TestToTransportMessage(t *testing.T) {
	svc, _ := setupProduceSvc()

	dp := &datapoint.Datapoint{
		FunctionGroup:  50,
		FunctionNumber: 0,
		DatapointID:    40650,
		Writable:       true,
		TypeName:       "LIST",
	}

	msg := hoval.NewMessage(
		1153,
		1,
		0x40,
		dp,
		1,
	)

	tMsg, err := svc.ToTransportMessage(msg)
	require.NoError(t, err)
	require.NotNil(t, tMsg)

	assert.Equal(t, uint32(0x1fe40801), tMsg.ID)
	assert.Equal(t, uint8(7), tMsg.Length)
	assert.Equal(t, transport.Data{0x1, 0x40, 0x32, 0x0, 0x9e, 0xca, 0x1, 0x0}, tMsg.Data)
}

func TestToTransportMessageInvalidData(t *testing.T) {
	svc, _ := setupProduceSvc()

	msg := hoval.NewMessage(
		1153,
		1,
		0x40,
		nil,
		12,
	)

	_, err := svc.ToTransportMessage(msg)
	require.Error(t, err)
}
