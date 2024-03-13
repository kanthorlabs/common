package config

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := Default(testdata.PostgresUri)
		require.NoError(st, conf.Validate())
	})

	t.Run("KO - connection error", func(st *testing.T) {
		conf := Default(testdata.PostgresUri)
		conf.Connection.MaxLifetime = 0
		require.Error(st, conf.Validate())
	})
}

func TestDefault(t *testing.T) {
	conf := Default(testdata.PostgresUri)
	require.Equal(t, testdata.PostgresUri, conf.Uri)
	require.Equal(t, DefaultConnMaxLifetime, conf.Connection.MaxLifetime)
	require.Equal(t, DefaultConnMaxIdletime, conf.Connection.MaxIdletime)
	require.Equal(t, DefaultConnMaxIdleCount, conf.Connection.MaxIdleCount)
	require.Equal(t, DefaultConnMaxOpenCount, conf.Connection.MaxOpenCount)
}
