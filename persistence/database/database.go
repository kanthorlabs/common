package database

import (
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/database/config"
	"github.com/kanthorlabs/common/persistence/sqlx"
)

func New(provider configuration.Provider) (Database, error) {
	conf, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}

	return sqlx.New(conf.Sqlx, logger)
}

type Database interface {
	persistence.Persistence
}
