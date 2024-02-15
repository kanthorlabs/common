package database

import (
	"testing"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/persistence/database/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		provider, err := configuration.New(testdata.Fake.Color().SafeColorName())
		require.Nil(st, err)

		provider.SetDefault("database.engine", config.EngineSqlx)
		provider.SetDefault("database.sqlx.uri", testdata.SqliteUri)
		provider.SetDefault("logger.level", "fatal")

		_, err = New(provider)
		require.Nil(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		provider, err := configuration.New(testdata.Fake.Color().SafeColorName())
		require.Nil(st, err)
		_, err = New(provider)
		require.ErrorContains(t, err, "DATABASE.CONFIG")
	})

	t.Run("KO - logger error", func(st *testing.T) {
		provider, err := configuration.New(testdata.Fake.Color().SafeColorName())
		require.Nil(st, err)

		provider.SetDefault("database.engine", config.EngineSqlx)
		provider.SetDefault("database.sqlx.uri", testdata.SqliteUri)

		_, err = New(provider)
		require.ErrorContains(t, err, "LOGGER.CONFIG")
	})
}
