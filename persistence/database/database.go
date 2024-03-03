package database

import (
	"errors"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/database/config"
	"github.com/kanthorlabs/common/persistence/sqlx"
)

func New(conf *config.Config, logger logging.Logger) (Database, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	if conf.Engine == config.EngineSqlx {
		return sqlx.New(conf.Sqlx, logger.With("database", "sqlx"))
	}

	return nil, errors.New("DATABASE.ENGINE_UNKNOWN.ERROR")
}

type Database interface {
	persistence.Persistence
}
