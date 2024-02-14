package config

import (
	"github.com/kanthorlabs/common/configuration"
	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/validator"
)

var (
	EngineSqlx = "sqlx"
)

func New(provider configuration.Provider) (*Config, error) {
	provider.SetDefault("database.sqlx.connection.max_lifetime", sqlx.DefaultConnMaxLifetime)
	provider.SetDefault("database.sqlx.connection.max_idletime", sqlx.DefaultConnMaxIdletime)
	provider.SetDefault("database.sqlx.connection.max_idle_count", sqlx.DefaultConnMaxIdleCount)
	provider.SetDefault("database.sqlx.connection.max_open_count", sqlx.DefaultConnMaxOpenCount)

	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &conf.Database, nil
}

type Wrapper struct {
	Database Config `json:"database" yaml:"database" mapstructure:"database"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Database.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Engine string       `json:"engine" yaml:"engine" mapstructure:"engine"`
	Sqlx   *sqlx.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringOneOf("DATABASE.CONFIG.ENGINE", conf.Engine, []string{EngineSqlx}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineSqlx {
		if err := conf.Sqlx.Validate(); err != nil {
			return err
		}
	}

	return nil
}
