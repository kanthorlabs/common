package sqlx

import (
	"testing"

	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestNewGorm(t *testing.T) {
	t.Run("OK - postgres", func(st *testing.T) {
		conf := config.Default(testdata.PostgresUri + "?skip_default_transaction=true")
		_, err := NewGorm(conf, testify.Logger())
		require.ErrorContains(st, err, "postgres")

		require.NotContains(st, conf.Uri, "skip_default_transaction")
	})

	t.Run("OK - memory", func(st *testing.T) {
		conf := config.Default(testdata.SqliteUri + "&skip_default_transaction=true")
		db, err := NewGorm(conf, testify.Logger())
		require.NoError(st, err)
		require.Equal(st, "sqlite", db.Dialector.Name())

		require.NotContains(st, conf.Uri, "skip_default_transaction")
	})
}
