package sqlx

import (
	"context"
	"testing"

	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSqlx_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := config.Default(testdata.SqliteUri)
		instance, err := New(conf, testify.Logger())
		require.NoError(st, err)
		require.NotNil(st, instance)
		require.Implements(st, (*persistence.Persistence)(nil), instance)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := New(&config.Config{}, testify.Logger())
		assert.NotNil(t, err)
	})
}

func TestSqlx_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		instance := start(t)
		defer end(t, instance)

		require.NoError(st, instance.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		instance := start(t)
		defer end(t, instance)

		require.NoError(st, instance.Connect(context.Background()))
		require.ErrorIs(st, instance.Connect(context.Background()), ErrAlreadyConnected)
	})

	t.Run("KO - connection error", func(st *testing.T) {
		conf := config.Default(testdata.PostgresUri)
		instance, err := New(conf, testify.Logger())
		require.NoError(st, err)
		require.NotNil(st, instance)

		require.ErrorContains(st, instance.Connect(context.Background()), "SQLX.CONNECT.ERROR")
	})
}

func TestSqlx_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		instance := start(t)
		require.NoError(st, instance.Connect(context.Background()))
		require.NoError(st, instance.Readiness())
		require.NoError(st, instance.Disconnect(context.Background()))
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		instance := start(t)
		require.NoError(st, instance.Connect(context.Background()))
		require.NoError(st, instance.Disconnect(context.Background()))
		require.NoError(st, instance.Readiness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		instance := start(t)
		require.ErrorIs(st, instance.Readiness(), ErrNotConnected)
	})
}

func TestSqlx_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		instance := start(t)
		require.NoError(st, instance.Connect(context.Background()))
		require.NoError(st, instance.Liveness())
		require.NoError(st, instance.Disconnect(context.Background()))
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		instance := start(t)
		require.NoError(st, instance.Connect(context.Background()))
		require.NoError(st, instance.Disconnect(context.Background()))
		require.NoError(st, instance.Liveness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		instance := start(t)
		require.ErrorIs(st, instance.Liveness(), ErrNotConnected)
	})
}

func TestSqlx_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		instance := start(t)
		require.NoError(st, instance.Connect(context.Background()))
		require.NoError(st, instance.Disconnect(context.Background()))
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		instance := start(t)
		require.ErrorIs(st, instance.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestSqlx_Client(t *testing.T) {
	instance := start(t)

	require.Nil(t, instance.Client())
	require.NoError(t, instance.Connect(context.Background()))
	require.NotNil(t, instance.Client())
	require.NoError(t, instance.Disconnect(context.Background()))
	require.Nil(t, instance.Client())
}

func start(t *testing.T) persistence.Persistence {
	conf := config.Default(testdata.SqliteUri)
	instance, err := New(conf, testify.Logger())
	require.NoError(t, err)
	return instance
}

func end(t *testing.T, instance persistence.Persistence) {
	conn, err := instance.Client().(*gorm.DB).DB()
	require.NoError(t, err)
	err = conn.Close()
	require.NoError(t, err)
}
