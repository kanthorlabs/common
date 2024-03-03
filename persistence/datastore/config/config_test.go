package config

import (
	"testing"

	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := Config{
			Engine: EngineSqlx,
			Sqlx:   sqlxconfig.Default(testdata.SqliteUri),
		}
		require.NoError(st, conf.Validate())
	})

	t.Run("KO - engine unknown error", func(st *testing.T) {
		conf := Config{}
		require.ErrorContains(st, conf.Validate(), "DATASTORE.CONFIG")
	})

	t.Run("KO - sqlx configuration error", func(st *testing.T) {
		conf := Config{
			Engine: EngineSqlx,
			Sqlx:   &sqlxconfig.Config{},
		}
		require.ErrorContains(st, conf.Validate(), "SQLX.CONFIG.URI")
	})
}
