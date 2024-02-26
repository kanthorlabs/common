package config

import "github.com/kanthorlabs/common/validator"

type Config struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.StringUri("CACHE.CONFIG.URI", conf.Uri),
	)
}
