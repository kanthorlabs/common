package config

import (
	"testing"

	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	sqlxconf := sqlxconfig.Default(testdata.SqliteUri)
	t.Run("OK", func(st *testing.T) {
		conf := Config{Sqlx: *sqlxconf}
		require.NoError(st, conf.Validate())
	})

	t.Run("KO - sqlx configuration error", func(st *testing.T) {
		conf := Config{
			Sqlx: sqlxconfig.Config{},
		}
		require.ErrorContains(st, conf.Validate(), "SQLX.CONFIG.URI")
	})
}
