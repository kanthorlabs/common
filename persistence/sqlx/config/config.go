package config

import "github.com/kanthorlabs/common/validator"

var (
	TypePostgres = "postgres"
	TypeSqlite   = "file"

	DefaultConnMaxLifetime  int64 = 300000
	DefaultConnMaxIdletime  int64 = 60000
	DefaultConnMaxIdleCount int   = 1
	DefaultConnMaxOpenCount int   = 10
)

type Config struct {
	Uri            string     `json:"uri" yaml:"uri" mapstructure:"uri"`
	SkipDefaultTxn bool       `json:"skip_default_txn" yaml:"skip_default_txn" mapstructure:"skip_default_txn"`
	Connection     Connection `json:"connection" yaml:"connection" mapstructure:"connection"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringUri("SQLX.CONFIG.URI", conf.Uri),
		validator.StringStartsWithOneOf("SQLX.CONFIG.URI", conf.Uri, []string{TypePostgres, TypeSqlite}),
	)
	if err != nil {
		return err
	}

	if err := conf.Connection.Validate(); err != nil {
		return err
	}

	return nil
}
