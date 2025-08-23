package config_test

import (
	"testing"

	"github.com/chrishrb/hoval-gateway/config"
	clone "github.com/huandu/go-clone/generic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigure(t *testing.T) {
	cfg := clone.Clone(&config.DefaultConfig)

	settings, err := config.Configure(t.Context(), cfg)
	require.NoError(t, err)

	wantApiSettings := config.HttpApiSettings{
		Addr:    "localhost:8080",
		Host:    "localhost",
		OrgName: "chrishrb",
	}

	assert.Equal(t, wantApiSettings, settings.HttpApi)
	assert.NotNil(t, settings.Storage)
	assert.NotNil(t, settings.DatapointProvider)
	assert.NotNil(t, settings.TransportConsumer)
	assert.NotNil(t, settings.TransportSender)
}

func TestConfigureInMemoryStorage(t *testing.T) {
	cfg := clone.Clone(&config.DefaultConfig)
	cfg.Storage.Type = "in_memory"

	settings, err := config.Configure(t.Context(), cfg)
	require.NoError(t, err)
	require.NotNil(t, settings.Storage)
}
