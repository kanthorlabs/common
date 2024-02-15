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

func TestSqlite(t *testing.T) {
	t.Run("Ko because of configuration validation", func(st *testing.T) {
		_, err := New(&config.Config{}, testify.Logger())
		assert.NotNil(t, err)
	})

	t.Run("readiness", func(st *testing.T) {
		instance := start(st)
		require.ErrorIs(st, instance.Readiness(), ErrNotConnected)

		// connect then check readiness
		require.Nil(st, instance.Connect(context.Background()))
		require.Nil(st, instance.Readiness())

		// disconnect then check readiness
		require.Nil(st, instance.Disconnect(context.Background()))
		require.Nil(st, instance.Readiness())

		// close the connection manually
		require.Nil(st, instance.Connect(context.Background()))
		end(st, instance)

		// the readiness should be failed
		require.ErrorIs(st, instance.Readiness(), ErrNotReady)
	})

	t.Run("liveness", func(st *testing.T) {
		instance := start(st)
		require.ErrorIs(st, instance.Liveness(), ErrNotConnected)

		// connect then check readiness
		require.Nil(st, instance.Connect(context.Background()))
		require.Nil(st, instance.Liveness())

		// disconnect then check readiness
		require.Nil(st, instance.Disconnect(context.Background()))
		require.Nil(st, instance.Liveness())

		// close the connection manually
		require.Nil(st, instance.Connect(context.Background()))
		end(st, instance)

		// the readiness should be failed
		require.ErrorIs(st, instance.Liveness(), ErrNotLive)

	})

	t.Run("connection", func(st *testing.T) {
		instance := start(st)
		ctx := context.Background()

		// unabel to disconnect if you didn't connect first
		require.ErrorIs(st, instance.Disconnect(ctx), ErrNotConnected)

		require.Nil(st, instance.Connect(ctx))
		// already connect, should not start new connection
		require.ErrorIs(st, instance.Connect(ctx), ErrAlreadyConnected)

		require.Nil(st, instance.Disconnect(ctx))
	})
}

func start(t *testing.T) persistence.Persistence {
	conf := &config.Config{
		Uri: testdata.SqliteUri,
		Connection: config.Connection{
			MaxLifetime:  config.DefaultConnMaxLifetime,
			MaxIdletime:  config.DefaultConnMaxIdletime,
			MaxIdleCount: config.DefaultConnMaxIdleCount,
			MaxOpenCount: config.DefaultConnMaxOpenCount,
		},
	}
	instance, err := New(conf, testify.Logger())
	require.Nil(t, err)

	return instance
}

func end(t *testing.T, instance persistence.Persistence) {
	conn, err := instance.Client().(*gorm.DB).DB()
	require.Nil(t, err)
	err = conn.Close()
	require.Nil(t, err)
}
