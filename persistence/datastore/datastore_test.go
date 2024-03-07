package datastore

import (
	"testing"

	"github.com/kanthorlabs/common/persistence/datastore/config"
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestDatastore_New(t *testing.T) {
	sqlxconf := sqlxconfig.Default(testdata.SqliteUri)
	t.Run("OK", func(st *testing.T) {
		conf := &config.Config{
			Engine: config.EngineSqlx,
			Sqlx:   *sqlxconf,
		}
		_, err := New(conf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := New(&config.Config{}, testify.Logger())
		require.ErrorContains(st, err, "DATASTORE.CONFIG")
	})
}
