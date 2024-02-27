package config

import "github.com/kanthorlabs/common/validator"

type Config struct {
	Uri        string `json:"uri" yaml:"uri" mapstructure:"uri"`
	TimeToLive int    `json:"time_to_live" yaml:"time_to_live" mapstructure:"time_to_live"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.StringUri("IDEMPOTENCY.CONFIG.URI", conf.Uri),
		validator.NumberGreaterThan("IDEMPOTENCY.CONFIG.TIME_TO_LIVE", conf.TimeToLive, 1000),
	)
}
