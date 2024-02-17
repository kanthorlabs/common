package config

import (
	"github.com/kanthorlabs/common/configuration"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &conf.Passport, nil
}

type Wrapper struct {
	Passport Config `json:"passport" yaml:"passport" mapstructure:"passport"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Passport.Validate(); err != nil {
		return err
	}
	return nil
}
