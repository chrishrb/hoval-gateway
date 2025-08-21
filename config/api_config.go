package config

type CorsConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins" json:"allowed_origins" validate:"required"`
	AllowedMethods []string `mapstructure:"allowed_methods" json:"allowed_methods" validate:"required"`
	AllowedHeaders []string `mapstructure:"allowed_headers" json:"allowed_headers" validate:"required"`
}

type HttpApiSettingsConfig struct {
	Addr    string      `mapstructure:"addr" json:"addr" validate:"required"`
	Host    string      `mapstructure:"host,omitempty" json:"host,omitempty"`
	OrgName string      `mapstructure:"org_name,omitempty" json:"org_name,omitempty"`
	Cors    *CorsConfig `mapstructure:"cors,omitempty" json:"cors,omitempty"`
}

type MqttSettingsConfig struct {
	Urls              []string `mapstructure:"urls" toml:"urls" validate:"required,dive,required"`
	Prefix            string   `mapstructure:"prefix" toml:"prefix" validate:"required"`
	Group             string   `mapstructure:"group" toml:"group" validate:"required"`
	ConnectTimeout    string   `mapstructure:"connect_timeout" toml:"connect_timeout" validate:"required"`
	ConnectRetryDelay string   `mapstructure:"connect_retry_delay" toml:"connect_retry_delay" validate:"required"`
	KeepAliveInterval string   `mapstructure:"keep_alive_interval" toml:"keep_alive_interval" validate:"required"`
}

type PubSubApiConfig struct {
	Type string              `mapstructure:"type" toml:"type" validate:"required,oneof=mqtt"`
	Mqtt *MqttSettingsConfig `mapstructure:"mqtt,omitempty" toml:"mqtt,omitempty" validate:"required_if=Type mqtt"`
}

type ApiSettingsConfig struct {
	Http   HttpApiSettingsConfig `mapstructure:"http_api,omitempty" json:"http_api" validate:"required"`
	PubSub *PubSubApiConfig      `mapstructure:"pubsub_api,omitempty" json:"pubsub_api,omitempty"`
}
