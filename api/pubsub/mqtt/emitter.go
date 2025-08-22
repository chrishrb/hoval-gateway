package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/chrishrb/hoval-gateway/api/pubsub"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type Emitter struct {
	sync.Mutex
	connectionDetails
	conn *autopaho.ConnectionManager
}

func NewEmitter(opts ...Opt[Emitter]) pubsub.Emitter {
	e := new(Emitter)
	for _, opt := range opts {
		opt(e)
	}
	ensureEmitterDefaults(e)
	return e
}

func (e *Emitter) Emit(ctx context.Context, functionGroup, functionNumber uint8, datapointID uint16, message *pubsub.Message) error {
	topic := fmt.Sprintf("%s/out/%d/%d/%d", e.mqttPrefix, functionGroup, functionNumber, datapointID)
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshalling response of datapoint %d/%d/%d: %v", functionGroup, functionNumber, datapointID, err)
	}

	err = e.ensureConnection(ctx)
	if err != nil {
		return fmt.Errorf("connecting to MQTT: %v", err)
	}

	_, err = e.conn.Publish(ctx, &paho.Publish{
		Topic:   topic,
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("publishing to %s: %v", topic, err)
	}
	return nil
}

func ensureEmitterDefaults(e *Emitter) {
	if e.mqttBrokerUrls == nil {
		u, err := url.Parse("mqtt://127.0.0.1:1883/")
		if err != nil {
			panic(err)
		}
		e.mqttBrokerUrls = []*url.URL{u}
	}
	if e.mqttPrefix == "" {
		e.mqttPrefix = "cs"
	}
	if e.mqttConnectTimeout == 0 {
		e.mqttConnectTimeout = 10 * time.Second
	}
	if e.mqttConnectRetryDelay == 0 {
		e.mqttConnectRetryDelay = 1 * time.Second
	}
	if e.mqttKeepAliveInterval == 0 {
		e.mqttKeepAliveInterval = 10
	}
}

func (e *Emitter) ensureConnection(ctx context.Context) error {
	e.Lock()
	defer e.Unlock()
	if e.conn == nil {
		conn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
			BrokerUrls:        e.mqttBrokerUrls,
			KeepAlive:         e.mqttKeepAliveInterval,
			ConnectRetryDelay: e.mqttConnectRetryDelay,
			ClientConfig: paho.ClientConfig{
				ClientID: fmt.Sprintf("%s-%s", "hoval-gw-emit", randSeq(5)),
			},
		})
		if err != nil {
			return err
		}
		e.conn = conn

		err = conn.AwaitConnection(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
