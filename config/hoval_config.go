package config

type PeriodicRequest struct {
	ReceiverMask   uint32 `mapstructure:"receiver_mask" json:"receiver_mask" validate:"required"`
	FunctionGroup  uint8  `mapstructure:"function_group" json:"function_group" validate:"required"`
	FunctionNumber uint8  `mapstructure:"function_number" json:"function_number" validate:"required"`
	DatapointID    uint16 `mapstructure:"datapoint_id" json:"datapoint_id" validate:"required"`
}

type HovalConfig struct {
	SenderID         uint32            `mapstructure:"sender_id" json:"sender_id" validate:"required"`
	PeriodicRequests []PeriodicRequest `mapstructure:"periodic_requests" json:"periodic_requests" validate:"dive"`
	PeriodicRunEvery string            `mapstructure:"run_every" json:"run_every" validate:"required,gt=0"`
}
