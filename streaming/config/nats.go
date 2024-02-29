package config

import "github.com/kanthorlabs/common/validator"

type Nats struct {
	Replicas int        `json:"replicas" yaml:"replicas" mapstructure:"replicas"`
	Limits   NatsLimits `json:"limits" yaml:"limits" mapstructure:"limits"`
}

func (conf *Nats) Validate() error {
	err := validator.Validate(
		validator.NumberGreaterThanOrEqual("STREAMING.CONFIG.NATS.REPLICAS", conf.Replicas, 0),
	)
	if err != nil {
		return err
	}

	if err := conf.Limits.Validate(); err != nil {
		return err
	}

	return nil
}

type NatsLimits struct {
	Bytes    int64 `json:"bytes" yaml:"bytes" mapstructure:"bytes"`
	MsgSize  int32 `json:"msg_size" yaml:"msg_size" mapstructure:"msg_size"`
	MsgCount int64 `json:"msg_count" yaml:"msg_count" mapstructure:"msg_count"`
	MsgAge   int64 `json:"msg_age" yaml:"msg_age" mapstructure:"msg_age"`
}

func (conf *NatsLimits) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThanOrEqual("STREAMING.CONFIG.NATS.LIMITS.BYTES", conf.Bytes, 1024),
		validator.NumberGreaterThanOrEqual("STREAMING.CONFIG.NATS.LIMITS.MSG_SIZE", conf.MsgSize, 1024),
		validator.NumberGreaterThanOrEqual("STREAMING.CONFIG.NATS.LIMITS.MSG_COUNT", conf.MsgCount, 1),
		validator.NumberGreaterThanOrEqual("STREAMING.CONFIG.NATS.LIMITS.AGE", conf.MsgAge, 1000),
	)
}
