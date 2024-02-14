package config

import "github.com/kanthorlabs/common/validator"

type Config struct {
	Strategies []Strategy `json:"strategies" yaml:"strategies" mapstructure:"strategies"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.SliceRequired("PASSPORT.CONFIG.STRATEGIES", conf.Strategies),
		validator.Slice(conf.Strategies, func(_ int, strategy *Strategy) error {
			return strategy.Validate()
		}),
	)
}
