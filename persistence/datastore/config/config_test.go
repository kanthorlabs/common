package config

import (
	"testing"

	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	provider, err := configuration.New("kanthor")
	require.Nil(t, err)

	_, err = New(provider)
	require.ErrorContains(t, err, "DATASTORE.CONFIG.ENGINE")

	provider.SetDefault("datastore.engine", EngineSqlx)

	_, err = New(provider)
	require.ErrorContains(t, err, "SQLX.CONFIG.URI")

	provider.SetDefault("datastore.sqlx.uri", testdata.SqliteUri)

	_, err = New(provider)
	require.Nil(t, err)
}
