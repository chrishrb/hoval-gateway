package config

type InMemoryStorageConfig struct{}

type StorageConfig struct {
	Type            string                 `mapstructure:"type" json:"type" validate:"required,oneof=in_memory"`
	InMemoryStorage *InMemoryStorageConfig `mapstructure:"in_memory,omitempty" json:"in_memory,omitempty"`
}
