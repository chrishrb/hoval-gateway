package mock

import (
	"github.com/chrishrb/hoval-gateway/transport"
)

type MockBus struct {
	Bus chan transport.Message
}

func NewMockBus() *MockBus {
	return &MockBus{
		Bus: make(chan transport.Message),
	}
}

func (m *MockBus) Close() error {
	close(m.Bus)
	return nil
}
