package datastore

import (
	"errors"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/datastore/config"
	"github.com/kanthorlabs/common/persistence/sqlx"
)

func New(conf *config.Config, logger logging.Logger) (Datastore, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	if conf.Engine == config.EngineSqlx {
		return sqlx.New(conf.Sqlx, logger.With("database", "sqlx"))
	}

	return nil, errors.New("DATASTORE.ENGINE_UNKNOWN.ERROR")
}

type Datastore interface {
	persistence.Persistence
}
