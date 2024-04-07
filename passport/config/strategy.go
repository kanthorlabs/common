package config

import (
	"github.com/kanthorlabs/common/validator"
)

var (
	EngineAsk      = "ask"
	EngineInternal = "internal"
	EngineExternal = "external"
)

type Strategy struct {
	Name   string `json:"name" yaml:"name" mapstructure:"name"`
	Engine string `json:"engine" yaml:"engine" mapstructure:"engine"`

	Ask      Ask      `json:"ask" yaml:"ask" mapstructure:"ask"`
	Internal Internal `json:"internal" yaml:"internal" mapstructure:"internal"`
	External External `json:"external" yaml:"external" mapstructure:"external"`
}

func (conf *Strategy) Validate() error {
	err := validator.Validate(
		validator.StringRequired("PASSPORT.STRATEGY.CONFIG.NAME", conf.Name),
		validator.StringOneOf("PASSPORT.STRATEGY.CONFIG.ENGINE", conf.Engine, []string{EngineAsk, EngineInternal, EngineExternal}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		if err := conf.Ask.Validate(); err != nil {
			return err
		}
	}

	if conf.Engine == EngineInternal {
		if err := conf.Internal.Validate(); err != nil {
			return err
		}
	}

	return nil
}
