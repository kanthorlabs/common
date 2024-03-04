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
		validator.StringRequired("GATEWAY.CONFIG.ADDR", conf.Addr),
		validator.NumberGreaterThanOrEqual("GATEWAY.CONFIG.TIMEOUT", conf.Timeout, 1000),
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
	AllowedOrigins   []string `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers" yaml:"allowed_headers"`
	ExposedHeaders   []string `json:"exposed_origins" yaml:"exposed_origins"`
	AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials"`
	MaxAge           int64    `json:"max_age" yaml:"max_age"`
}

func (conf *Cors) Validate() error {
	return validator.Validate(
		validator.NumberInRange("GATEWAY.CONFIG.CORS.MAX_AGE", conf.MaxAge, 1000, 86400000),
	)
}
