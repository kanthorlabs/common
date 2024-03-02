package datastore

import (
	"testing"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/persistence/datastore/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestDatastore(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		provider, err := configuration.New(testdata.Fake.Color().SafeColorName())
		require.NoError(st, err)

		provider.SetDefault("datastore.engine", config.EngineSqlx)
		provider.SetDefault("datastore.sqlx.uri", testdata.SqliteUri)
		provider.SetDefault("logger.level", "fatal")

		_, err = New(provider)
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		provider, err := configuration.New(testdata.Fake.Color().SafeColorName())
		require.NoError(st, err)
		_, err = New(provider)
		require.ErrorContains(t, err, "DATASTORE.CONFIG")
	})

	t.Run("KO - logger error", func(st *testing.T) {
		provider, err := configuration.New(testdata.Fake.Color().SafeColorName())
		require.NoError(st, err)

		provider.SetDefault("datastore.engine", config.EngineSqlx)
		provider.SetDefault("datastore.sqlx.uri", testdata.SqliteUri)

		_, err = New(provider)
		require.ErrorContains(t, err, "LOGGER.CONFIG")
	})
}
