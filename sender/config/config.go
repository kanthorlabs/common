package config

import (
	"github.com/kanthorlabs/common/validator"
)

type Config struct {
	Timeout int64             `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Headers map[string]string `json:"header" yaml:"header" mapstructure:"header"`
	Retry   Retry             `json:"retry" yaml:"retry" mapstructure:"retry"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.NumberGreaterThanOrEqual("SENDER.CONFIG.TIMEOUT", conf.Timeout, 1000),
	)
	if err != nil {
		return err
	}

	if err := conf.Retry.Validate(); err != nil {
		return err
	}

	return nil
}
