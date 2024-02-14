package config

import (
	"github.com/kanthorlabs/common/validator"
)

var (
	EngineAsk = "ask"
)

type Strategy struct {
	Engine string `json:"engine" yaml:"engine" mapstructure:"engine"`
	Name   string `json:"name" yaml:"name" mapstructure:"name"`

	Ask Ask `json:"ask" yaml:"ask" mapstructure:"ask"`
}

func (conf *Strategy) Validate() error {
	err := validator.Validate(
		validator.StringOneOf("PASSPORT.STRATEGY.CONFIG.ENGINE", conf.Engine, []string{EngineAsk}),
		validator.StringRequired("PASSPORT.STRATEGY.CONFIG.NAME", conf.Name),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		if err := conf.Ask.Validate(); err != nil {
			return err
		}
	}

	return nil
}
