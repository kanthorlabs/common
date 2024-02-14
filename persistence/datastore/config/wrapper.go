package config

import "github.com/kanthorlabs/common/configuration"

type Wrapper struct {
	Datastore Config `json:"datastore" yaml:"datastore" mapstructure:"datastore"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Datastore.Validate(); err != nil {
		return err
	}
	return nil
}

func New(provider configuration.Provider) (*Config, error) {
	// you will gain about 30%+ performance improvement after that by disable default txn
	provider.SetDefault("datastore.sqlx.skip_default_txn", true)

	provider.SetDefault("datastore.sqlx.connection.max_lifetime", 300000)
	provider.SetDefault("datastore.sqlx.connection.max_idletime", 60000)
	provider.SetDefault("datastore.sqlx.connection.max_idle_count", 1)
	provider.SetDefault("datastore.sqlx.connection.max_open_count", 10)

	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &conf.Datastore, nil
}
