package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	transport "github.com/chrishrb/hoval-gateway/api/pubsub"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type Listener struct {
	connectionDetails
	mqttGroup string
}

func NewListener(opts ...Opt[Listener]) *Listener {
	l := new(Listener)
	for _, opt := range opts {
		opt(l)
	}
	ensureListenerDefaults(l)
	return l
}

func ensureListenerDefaults(l *Listener) {
	if l.mqttBrokerUrls == nil {
		u, err := url.Parse("mqtt://127.0.0.1:1883/")
		if err != nil {
			panic(err)
		}
		l.mqttBrokerUrls = []*url.URL{u}
	}
	if l.mqttPrefix == "" {
		l.mqttPrefix = "hoval"
	}
	if l.mqttGroup == "" {
		l.mqttGroup = "hoval-gw"
	}
	if l.mqttConnectTimeout == 0 {
		l.mqttConnectTimeout = 10 * time.Second
	}
	if l.mqttConnectRetryDelay == 0 {
		l.mqttConnectRetryDelay = 1 * time.Second
	}
	if l.mqttKeepAliveInterval == 0 {
		l.mqttKeepAliveInterval = 10
	}
}

func (l *Listener) Connect(ctx context.Context, handler transport.MessageHandler) (transport.Connection, error) {
	var err error

	ctx, cancel := context.WithTimeout(ctx, l.mqttConnectTimeout)
	defer cancel()

	clientId := fmt.Sprintf("%s-%s", l.mqttGroup, randSeq(5))

	readyCh := make(chan struct{})

	topic := fmt.Sprintf("%s/in/#", l.mqttPrefix)

	conn := new(connection)
	mqttRouter := paho.NewStandardRouter()
	conn.mqttConn, err = autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        l.mqttBrokerUrls,
		KeepAlive:         l.mqttKeepAliveInterval,
		ConnectRetryDelay: l.mqttConnectRetryDelay,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err := manager.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{{Topic: topic}},
			})
			if err != nil {
				slog.Error("failed to subscribe to topic", "topic", topic)
				return
			}
			mqttRouter.UnregisterHandler(topic)
			mqttRouter.RegisterHandler(topic, func(mqttMsg *paho.Publish) {
				ctx := context.Background()

				// determine functionGroup, functionNumber, datapointID
				topicParts := strings.Split(mqttMsg.Topic, "/")

				functionGroup, err := stringToUint8((topicParts[len(topicParts)-3]))
				if err != nil {
					slog.Error("unable to convert function group to uint8", "err", err)
					return
				}
				functionNumber, err := stringToUint8(topicParts[len(topicParts)-2])
				if err != nil {
					slog.Error("unable to convert function number to uint8", "err", err)
					return
				}
				datapointID, err := stringToUint16(topicParts[len(topicParts)-1])
				if err != nil {
					slog.Error("unable to convert datapointID to uint16", "err", err)
					return
				}

				// unmarshal the message
				var msg transport.Message
				err = json.Unmarshal(mqttMsg.Payload, &msg)
				if err != nil {
					slog.Warn("unable to unmarshal message", "err", err)
					return
				}

				// execute the handler
				handler.Handle(ctx, functionGroup, functionNumber, datapointID, &msg)
			})
			readyCh <- struct{}{}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: clientId,
			Router:   mqttRouter,
		},
	})
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, errors.New("timeout waiting for mqtt connectionDetails setup")
	case <-readyCh:
		return conn, nil
	}
}

type connection struct {
	mqttConn *autopaho.ConnectionManager
}

func (c *connection) Disconnect(ctx context.Context) error {
	if c.mqttConn != nil {
		err := c.mqttConn.Disconnect(ctx)
		if err != nil {
			return err
		}
		c.mqttConn = nil
	}
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		//#nosec G404 - client suffix does not require secure random number generator
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func stringToUint8(in string) (uint8, error) {
	if in == "" {
		return 0, nil
	}
	i, err := strconv.ParseUint(in, 0, 8)
	if err != nil {
		return 0, err
	}
	return uint8(i), nil
}

func stringToUint16(in string) (uint16, error) {
	if in == "" {
		return 0, nil
	}
	i, err := strconv.ParseUint(in, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint16(i), nil
}
