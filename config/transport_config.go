package config

type CANSettingsConfig struct {
	Device string `mapstructure:"device" json:"device" validate:"required"`
}

type TransportConfig struct {
	Type string             `mapstructure:"type" json:"type" validate:"required,oneof=can mock"`
	CAN  *CANSettingsConfig `mapstructure:"can" json:"can"`
}
