package config

import (
	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/validator"
)

type Config struct {
	Engine string       `json:"engine" yaml:"engine" mapstructure:"engine"`
	Sqlx   *sqlx.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringOneOf("DATABASE.CONFIG.ENGINE", conf.Engine, []string{sqlx.Engine}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == sqlx.Engine {
		if err := conf.Sqlx.Validate(); err != nil {
			return err
		}
	}

	return nil
}
