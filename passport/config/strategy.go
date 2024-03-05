package config

import (
	"github.com/kanthorlabs/common/validator"
)

var (
	EngineAsk        = "ask"
	EngineDurability = "durability"
)

type Strategy struct {
	Name   string `json:"name" yaml:"name" mapstructure:"name"`
	Engine string `json:"engine" yaml:"engine" mapstructure:"engine"`

	Ask        Ask        `json:"ask" yaml:"ask" mapstructure:"ask"`
	Durability Durability `json:"durability" yaml:"durability" mapstructure:"durability"`
}

func (conf *Strategy) Validate() error {
	err := validator.Validate(
		validator.StringRequired("PASSPORT.STRATEGY.CONFIG.NAME", conf.Name),
		validator.StringOneOf("PASSPORT.STRATEGY.CONFIG.ENGINE", conf.Engine, []string{EngineAsk, EngineDurability}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		if err := conf.Ask.Validate(); err != nil {
			return err
		}
	}

	if conf.Engine == EngineDurability {
		if err := conf.Durability.Validate(); err != nil {
			return err
		}
	}

	return nil
}
