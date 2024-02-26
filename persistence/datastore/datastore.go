package datastore

import (
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/datastore/config"
	"github.com/kanthorlabs/common/persistence/sqlx"
)

func New(provider configuration.Provider) (Datastore, error) {
	conf, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	logger, err := logging.New(provider)
	if err != nil {
		return nil, err
	}

	return sqlx.New(conf.Sqlx, logger.With("datastore", "sqlx"))
}

type Datastore interface {
	persistence.Persistence
}
