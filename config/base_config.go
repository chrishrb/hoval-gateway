package config

import (
	"bufio"
	"io"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert/yaml"
)

type BaseConfig struct {
	Api       ApiSettingsConfig `mapstructure:"api" json:"api" validate:"required"`
	Transport TransportConfig   `mapstructure:"transport" json:"transport" validate:"required"`
	Storage   StorageConfig     `mapstructure:"storage" json:"storage" validate:"required"`
	Hoval     HovalConfig       `mapstructure:"hoval" json:"hoval" validate:"required"`
}

// DefaultConfig provides the default configuration. The configuration
// read from the YAML file will overlay this configuration.
var DefaultConfig = BaseConfig{
	Api: ApiSettingsConfig{
		Http: HttpApiSettingsConfig{
			Addr:    "localhost:8080",
			Host:    "localhost",
			OrgName: "chrishrb",
		},
	},
	Transport: TransportConfig{
		Type: "mock",
	},
	Storage: StorageConfig{
		Type: "in_memory",
	},
	Hoval: HovalConfig{
		SenderID: 1153, // Default for Hoval GW
	},
}

// Load reads YAML configuration from a reader.
func (c *BaseConfig) Load(reader io.Reader) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, c); err != nil {
		return err
	}
	return nil
}

// LoadFromFile reads YAML configuration from a file.
func (c *BaseConfig) LoadFromFile(configFile string) error {
	//#nosec G304 - only files specified by the person running the application will be loaded
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}
	err = c.Load(bufio.NewReader(f))
	return err
}

// Validate ensures that the configuration is structurally valid.
func (c *BaseConfig) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}
