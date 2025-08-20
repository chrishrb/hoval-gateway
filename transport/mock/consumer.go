package mock

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/chrishrb/hoval-gateway/transport"
)

type Consumer struct {
	conn *MockBus
}

func NewConsumer(bus *MockBus) *Consumer {
	return &Consumer{
		conn: bus,
	}
}

func (c *Consumer) Consume(ctx context.Context, handler transport.MessageHandler) (transport.Connection, error) {
	go func() {
		for {
			msg := <-c.conn.Bus
			handler.Handle(ctx, &msg)
		}
	}()

	select {
	case <-ctx.Done():
		return nil, errors.New("timeout waiting for mock setup")
	default:
		return c.conn, nil
	}
}

func (c *Consumer) ReadFromFile(ctx context.Context, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l, err := parseCandumpLine(scanner.Text())
		if err != nil {
			return fmt.Errorf("error parsing line: %w", err)
		}
		c.conn.Bus <- *l
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func parseCandumpLine(line string) (*transport.Message, error) {
	candumpRegex := regexp.MustCompile(`\(([0-9.]+)\)\s+\w+\s+([0-9A-F]+)#([0-9A-F]*)`)

	matches := candumpRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("invalid format: %s", line)
	}

	// Parse CAN ID
	canID, err := strconv.ParseUint(matches[2], 16, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid CAN ID: %v", err)
	}

	// Parse data
	data, err := hex.DecodeString(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid data: %v", err)
	}

	s := make([]byte, 8)
	copy(s, data)

	return &transport.Message{
		ID:     uint32(canID),
		Length: uint8(len(data)),
		Data:   transport.Data(s),
	}, nil
}
