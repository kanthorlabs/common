package config

import (
	"github.com/kanthorlabs/common/validator"
)

type Config struct {
	Addr    string            `json:"addr" yaml:"addr"`
	Timeout int64             `json:"timeout" yaml:"timeout"`
	Headers map[string]string `json:"header" yaml:"header" mapstructure:"header"`
	Retry   Retry             `json:"retry" yaml:"retry"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringRequired("SENDER.CONFIG.ADDR", conf.Addr),
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
