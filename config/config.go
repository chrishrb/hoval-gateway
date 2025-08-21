package config

import (
	"context"
	"fmt"

	"github.com/chrishrb/hoval-gateway/store"
	"github.com/chrishrb/hoval-gateway/store/inmemory"
	"github.com/chrishrb/hoval-gateway/transport"
	"github.com/chrishrb/hoval-gateway/transport/can"
	"github.com/chrishrb/hoval-gateway/transport/mock"
	"k8s.io/utils/clock"
)

type HttpApiSettings struct {
	Addr    string
	Host    string
	OrgName string
	Cors    *CorsConfig
}

type Config struct {
	HttpApi           HttpApiSettings
	Storage           store.Engine
	DatapointProvider DatapointProvider
	CanConsumer       transport.Consumer
	CanSender         transport.Sender
	// TODO: mqtt
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

	c.DatapointProvider, err = NewCsvDatapointProvider("datapoints.csv")
	if err != nil {
		return nil, err
	}

	c.CanConsumer, err = getCanConsumer(&cfg.Transport)
	if err != nil {
		return nil, err
	}

	c.CanSender, err = getCanSender(&cfg.Transport)
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
