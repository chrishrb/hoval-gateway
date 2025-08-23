package mqtt_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/chrishrb/hoval-gateway/api/pubsub"
	"github.com/chrishrb/hoval-gateway/api/pubsub/mqtt"
	server "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListenerProcessesMessagesReceivedFromTheBroker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// start the broker
	broker, clientUrl := mqtt.NewBroker(t)
	defer func() {
		err := broker.Close()
		assert.NoError(t, err)
	}()
	err := broker.Serve()
	require.NoError(t, err)

	// setup the handler
	receivedMsgCh := make(chan struct{})
	handler := func(ctx context.Context, receiverMask uint32, msg *pubsub.Message) {
		assert.Equal(t, uint32(3), receiverMask)
		assert.Equal(t, uint8(50), msg.FunctionGroup)
		assert.Equal(t, uint8(2), msg.FunctionNumber)
		assert.Equal(t, uint16(123), msg.DatapointID)
		assert.Equal(t, 123.23, msg.Data)
		receivedMsgCh <- struct{}{}
	}

	// connect the listener to the broker
	listener := mqtt.NewListener(mqtt.WithMqttBrokerUrl[mqtt.Listener](clientUrl))
	conn, err := listener.Connect(ctx, pubsub.MessageHandlerFunc(handler))
	require.NoError(t, err)
	defer func() {
		if conn != nil {
			err := conn.Disconnect(ctx)
			require.NoError(t, err)
		}
	}()

	// publish message
	publishMessage(t, broker, pubsub.Message{
		FunctionGroup:  50,
		FunctionNumber: 2,
		DatapointID:    123,
		Data:           123.23,
	})

	// wait for message to be received / timeout
	select {
	case <-ctx.Done():
		assert.Fail(t, "timeout waiting for test to complete")
	case <-receivedMsgCh:
		// do nothing
	}
}

func publishMessage(t *testing.T, broker *server.Server, msg pubsub.Message) {
	msgBytes, err := json.Marshal(msg)
	require.NoError(t, err)

	cl := broker.NewClient(nil, "local", "inline", true)
	err = broker.InjectPacket(cl, packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type:   packets.Publish,
			Qos:    0,
			Retain: false,
		},
		TopicName: "hoval/in/3",
		Payload:   msgBytes,
		PacketID:  uint16(0),
	})
	require.NoError(t, err)
}
