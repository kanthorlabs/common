package config

import (
	"github.com/kanthorlabs/common/validator"
)

type Config struct {
	Addr    string `json:"addr" yaml:"addr"`
	Timeout int64  `json:"timeout" yaml:"timeout"`
	Cors    Cors   `json:"cors" yaml:"cors"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringUri("GATEWAY.ADDR", conf.Addr),
		validator.NumberGreaterThanOrEqual("GATEWAY.TIMEOUT", conf.Timeout, 1000),
	)
	if err != nil {
		return err
	}

	if err := conf.Cors.Validate(); err != nil {
		return err
	}

	return nil
}

type Cors struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func (conf *Cors) Validate() error {
	return validator.Validate(
		validator.NumberInRange("GATEWAY.ADDR", conf.MaxAge, 1, 86400),
	)
}
