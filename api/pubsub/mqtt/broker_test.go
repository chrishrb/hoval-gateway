package mqtt_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/chrishrb/hoval-gateway/api/pubsub/mqtt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func TestNewBroker(t *testing.T) {
	broker, addr := mqtt.NewBroker(t)
	defer func() {
		err := broker.Close()
		if err != nil {
			t.Errorf("closing broker: %v", err)
		}
	}()

	err := broker.Serve()
	if err != nil {
		t.Fatalf("starting broker: %v", err)
	}

	doneCh := make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	router := paho.NewStandardRouter()
	router.RegisterHandler("test/#", func(m *paho.Publish) {
		t.Logf("received message: %v", m)
		doneCh <- struct{}{}
	})

	_, err = autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{addr},
		KeepAlive:         10,
		ConnectRetryDelay: 10,
		OnConnectionUp: func(manager *autopaho.ConnectionManager, connack *paho.Connack) {
			_, err = manager.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{{Topic: "test/#"}},
			})
			if err != nil {
				t.Fatalf("subscribing to messages: %v", err)
			}

			err = broker.Publish("test/123", []byte("test data"), false, 0)
			if err != nil {
				t.Errorf("publising message: %v", err)
			}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: "hv-1",
			Router:   router,
		},
	})
	if err != nil {
		t.Fatalf("connecting to broker: %v", err)
	}

	select {
	case <-doneCh:
		return
	case <-ctx.Done():
		t.Fatal("failed to receive message")
	}
}
