package config

import (
	"context"
	"fmt"
	"net/url"
	"time"
	_ "embed"

	"github.com/chrishrb/hoval-gateway/api/pubsub"
	mqtt2 "github.com/chrishrb/hoval-gateway/api/pubsub/mqtt"
	"github.com/chrishrb/hoval-gateway/hoval/datapoint"
	"github.com/chrishrb/hoval-gateway/hoval/service"
	"github.com/chrishrb/hoval-gateway/periodic"
	"github.com/chrishrb/hoval-gateway/store"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/chrishrb/hoval-gateway/transport/can"
	"github.com/chrishrb/hoval-gateway/transport/mock"
	"k8s.io/utils/clock"
)

//go:embed datapoints.csv
var datapointsFile []byte

type HttpApiSettings struct {
	Addr    string
	Host    string
	OrgName string
	Cors    *CorsConfig
}

type Config struct {
	HttpApi           HttpApiSettings
	Storage           store.Engine
	DatapointProvider datapoint.DatapointProvider
	TransportConsumer transport.Consumer
	TransportSender   transport.Sender
	ConsumeHandler    transport.MessageHandler
	PubSubListener    pubsub.Listener
	PubSubEmitter     pubsub.Emitter
	SendHandler       pubsub.MessageHandler
	PeriodicRequester *periodic.PeriodicRequester
}

func Configure(ctx context.Context, cfg *BaseConfig) (c *Config, err error) {
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	c = &Config{
		HttpApi: HttpApiSettings{
			Addr:    cfg.Api.Http.Addr,
			Host:    cfg.Api.Http.Host,
			OrgName: cfg.Api.Http.OrgName,
			Cors:    cfg.Api.Http.Cors,
		},
	}

	c.Storage, err = getStorage(ctx, &cfg.Storage)
	if err != nil {
		return nil, err
	}

	c.DatapointProvider, err = datapoint.NewCsvDatapointProviderFromBytes(datapointsFile)
	if err != nil {
		return nil, err
	}

	c.TransportConsumer, err = getCanConsumer(&cfg.Transport)
	if err != nil {
		return nil, err
	}

	c.TransportSender, err = getCanSender(&cfg.Transport)
	if err != nil {
		return nil, err
	}

	if cfg.Api.PubSub != nil {
		c.PubSubListener, err = getMqttReceiver(cfg.Api.PubSub)
		if err != nil {
			return nil, err
		}

		c.PubSubEmitter, err = getOcppMsgEmitter(cfg.Api.PubSub)
		if err != nil {
			return nil, err
		}
	}

	sendSvc := service.NewSendService(cfg.Hoval.SenderID, c.Storage, c.DatapointProvider, c.TransportSender)
	c.SendHandler = sendSvc

	c.PeriodicRequester, err = getPeriodicRequester(sendSvc, &cfg.Hoval)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getStorage(_ context.Context, cfg *StorageConfig) (engine store.Engine, err error) {
	switch cfg.Type {
	case "in_memory":
		engine = inmemory.NewStore(clock.RealClock{})
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Type)
	}

	return
}

func getCanConsumer(cfg *TransportConfig) (transport.Consumer, error) {
	switch cfg.Type {
	case "mock":
		mockBus := mock.NewMockBus()
		return mock.NewConsumer(mockBus), nil
	case "can":
		return can.NewConsumer(
			can.WithCANDevice[can.Consumer](cfg.CAN.Device),
		), nil
	default:
		return nil, fmt.Errorf("unknown transport type: %s", cfg.Type)
	}
}

func getCanSender(cfg *TransportConfig) (transport.Sender, error) {
	switch cfg.Type {
	case "mock":
		mockBus := mock.NewMockBus()
		return mock.NewSender(mockBus), nil
	case "can":
		return can.NewSender(
			can.WithCANDevice[can.Sender](cfg.CAN.Device),
		), nil
	default:
		return nil, fmt.Errorf("unknown transport type: %s", cfg.Type)
	}
}

func getMqttReceiver(cfg *PubSubApiConfig) (pubsub.Listener, error) {
	switch cfg.Type {
	case "mqtt":
		var mqttUrls []*url.URL
		for _, urlStr := range cfg.Mqtt.Urls {
			u, err := url.Parse(urlStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mqtt url: %w", err)
			}
			mqttUrls = append(mqttUrls, u)
		}

		mqttConnectTimeout, err := time.ParseDuration(cfg.Mqtt.ConnectTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect timeout: %w", err)
		}

		mqttConnectRetryDelay, err := time.ParseDuration(cfg.Mqtt.ConnectRetryDelay)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect retry delay: %w", err)
		}

		mqttKeepAliveInterval, err := time.ParseDuration(cfg.Mqtt.KeepAliveInterval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt keep alive interval: %w", err)
		}

		opts := []mqtt2.Opt[mqtt2.Listener]{
			mqtt2.WithMqttBrokerUrls[mqtt2.Listener](mqttUrls),
			mqtt2.WithMqttPrefix[mqtt2.Listener](cfg.Mqtt.Prefix),
			mqtt2.WithMqttConnectSettings[mqtt2.Listener](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval),
			mqtt2.WithMqttGroup(cfg.Mqtt.Group),
		}

		return mqtt2.NewListener(opts...), nil
	default:
		return nil, fmt.Errorf("unknown transport type: %s", cfg.Type)
	}
}
func getOcppMsgEmitter(cfg *PubSubApiConfig) (pubsub.Emitter, error) {
	switch cfg.Type {
	case "mqtt":
		var mqttUrls []*url.URL
		for _, urlStr := range cfg.Mqtt.Urls {
			u, err := url.Parse(urlStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse mqtt url: %w", err)
			}
			mqttUrls = append(mqttUrls, u)
		}

		mqttConnectTimeout, err := time.ParseDuration(cfg.Mqtt.ConnectTimeout)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect timeout: %w", err)
		}

		mqttConnectRetryDelay, err := time.ParseDuration(cfg.Mqtt.ConnectRetryDelay)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt connect retry delay: %w", err)
		}

		mqttKeepAliveInterval, err := time.ParseDuration(cfg.Mqtt.KeepAliveInterval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mqtt keep alive interval: %w", err)
		}

		mqttEmitter := mqtt2.NewEmitter(
			mqtt2.WithMqttBrokerUrls[mqtt2.Emitter](mqttUrls),
			mqtt2.WithMqttPrefix[mqtt2.Emitter](cfg.Mqtt.Prefix),
			mqtt2.WithMqttConnectSettings[mqtt2.Emitter](mqttConnectTimeout, mqttConnectRetryDelay, mqttKeepAliveInterval))

		return mqttEmitter, nil
	default:
		return nil, fmt.Errorf("unknown transport type: %s", cfg.Type)
	}
}

func getPeriodicRequester(svc *service.SendService, cfg *HovalConfig) (*periodic.PeriodicRequester, error) {
	runEvery, err := time.ParseDuration(cfg.PeriodicRunEvery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse periodic RunEvery: %w", err)
	}

	requests := make([]periodic.PeriodicRequest, len(cfg.PeriodicRequests))
	for i, r := range cfg.PeriodicRequests {
		requests[i] = periodic.PeriodicRequest{
			ReceiverMask:   r.ReceiverMask,
			FunctionGroup:  r.FunctionGroup,
			FunctionNumber: r.FunctionNumber,
			DatapointID:    r.DatapointID,
		}
	}

	return periodic.NewPeriodicRequester(svc, requests, runEvery), nil
}
