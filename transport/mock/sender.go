package mock

import (
	"context"

	"github.com/chrishrb/hoval-gateway/transport"
)

type Sender struct {
	bus *MockBus
}

func NewSender(bus *MockBus) *Sender {
	return &Sender{
		bus: bus,
	}
}

func (s *Sender) Send(ctx context.Context, message *transport.Message) error {
	s.bus.Bus <- *message
	return nil
}
