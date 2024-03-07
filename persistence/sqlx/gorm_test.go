package sqlx

import (
	"testing"

	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestGorm(t *testing.T) {
	t.Run("OK - postgres", func(st *testing.T) {
		conf := config.Default(testdata.PostgresUri)
		_, err := Gorm(conf, testify.Logger())
		require.ErrorContains(st, err, "postgres")
	})

	t.Run("OK - memory", func(st *testing.T) {
		conf := config.Default(testdata.SqliteUri)
		db, err := Gorm(conf, testify.Logger())
		require.NoError(st, err)
		require.Equal(st, "sqlite", db.Dialector.Name())
	})
}
