package config

import (
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
)

type Durability struct {
	Sqlx sqlxconfig.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Durability) Validate() error {
	return conf.Sqlx.Validate()
}
