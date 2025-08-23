package config

type HovalConfig struct {
	SenderID uint32 `mapstructure:"sender_id" json:"sender_id" validate:"required"`
}
