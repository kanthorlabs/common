package config

import (
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
)

type Config struct {
	Sqlx sqlxconfig.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Config) Validate() error {
	if err := conf.Sqlx.Validate(); err != nil {
		return err
	}

	return nil
}
