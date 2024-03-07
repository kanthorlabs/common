package database

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/database/config"
	"github.com/kanthorlabs/common/persistence/sqlx"
)

// New creates a new Database instance that is backed by the SQLX persistence layer.
// The database is mainly designed to work with well design structure data so most of time we should use SQL database, such as PostgreSQL, MySQL, etc.
func New(conf *config.Config, logger logging.Logger) (Database, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return sqlx.New(&conf.Sqlx, logger.With("database", "sqlx"))
}

type Database interface {
	persistence.Persistence
}
