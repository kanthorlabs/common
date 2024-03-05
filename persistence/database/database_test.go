package database

import (
	"testing"

	"github.com/kanthorlabs/common/persistence/database/config"
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestDatabase_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &config.Config{
			Engine: sqlxconfig.Engine,
			Sqlx:   sqlxconfig.Default(testdata.SqliteUri),
		}
		_, err := New(conf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := New(&config.Config{}, testify.Logger())
		require.ErrorContains(st, err, "DATABASE.CONFIG")
	})
}
