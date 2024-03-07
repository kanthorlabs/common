package datastore

import (
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/datastore/config"
	"github.com/kanthorlabs/common/persistence/sqlx"
)

// New creates a new Datastore instance that allow you to interact with the datastore layer
// The datastore is different from the database in that it is designed to work with different types of databases, not just SQL databases.
// For instance, some of the supported databases are: PostgreSQL, CockroachDB, Cassandra, etc.
// The reason is datastore often deal with high volume of data and need to be able to scale horizontally
// and write is mostly higher than read, so it is important to have a database that can handle high write throughput.
func New(conf *config.Config, logger logging.Logger) (Datastore, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return sqlx.New(&conf.Sqlx, logger.With("database", "sqlx"))
}

type Datastore interface {
	persistence.Persistence
}
