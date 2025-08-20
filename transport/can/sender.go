package can

import (
	"context"
	"net"
	"sync"

	"github.com/chrishrb/hoval-gateway/transport"
	go_can "go.einride.tech/can"
	"go.einride.tech/can/pkg/socketcan"
)

type Sender struct {
	sync.Mutex
	connectionDetails
	conn net.Conn
}

func NewSender(opts ...Opt[Sender]) *Sender {
	c := &Sender{}
	for _, opt := range opts {
		opt(c)
	}
	ensureSenderDefaults(c)
	return c
}

func ensureSenderDefaults(c *Sender) {
	if c.device == "" {
		c.device = "can0"
	}
}

func (s *Sender) Send(ctx context.Context, message *transport.Message) error {
	if err := s.ensureConnection(ctx); err != nil {
		return err
	}
	tx := socketcan.NewTransmitter(s.conn)
	frame := go_can.Frame{
		ID:         message.ID,
		Length:     message.Length,
		Data:       go_can.Data(message.Data),
		IsRemote:   false,
		IsExtended: true,
	}
	return tx.TransmitFrame(ctx, frame)
}

func (s *Sender) ensureConnection(ctx context.Context) error {
	s.Lock()
	defer s.Unlock()
	if s.conn == nil {
		conn, err := socketcan.DialContext(ctx, "can", s.device)
		if err != nil {
			return err
		}
		s.conn = conn
	}
	return nil
}
