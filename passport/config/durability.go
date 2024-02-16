package config

import (
	"github.com/kanthorlabs/common/persistence/sqlx/config"
)

type Durability struct {
	Sqlx config.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Durability) Validate() error {
	return conf.Sqlx.Validate()
}
