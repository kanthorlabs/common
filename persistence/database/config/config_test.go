package config

import (
	"testing"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	provider, err := configuration.New("kanthor")
	require.NoError(t, err)

	_, err = New(provider)
	require.ErrorContains(t, err, "DATABASE.CONFIG.ENGINE")

	provider.SetDefault("database.engine", EngineSqlx)

	_, err = New(provider)
	require.ErrorContains(t, err, "SQLX.CONFIG.")

	provider.SetDefault("database.sqlx.uri", testdata.SqliteUri)

	_, err = New(provider)
	require.NoError(t, err)
}
