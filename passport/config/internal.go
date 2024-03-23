package config

import (
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
)

type Internal struct {
	Sqlx sqlxconfig.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Internal) Validate() error {
	return conf.Sqlx.Validate()
}
