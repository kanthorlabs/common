package idempotency

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestRedis_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewRedis(conf, testify.Logger())
		require.ErrorContains(t, err, "IDEMPOTENCY.CONFIG.")
	})
}

func TestRedis_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		require.ErrorIs(t, c.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestRedis_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Readiness(), ErrNotConnected)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Disconnect(context.Background()))
		require.NoError(t, c.Readiness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Readiness(), ErrNotConnected)
	})
}

func TestRedis_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Liveness(), ErrNotConnected)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Disconnect(context.Background()))
		require.NoError(t, c.Liveness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Liveness(), ErrNotConnected)
	})
}

func TestRedis_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Disconnect(context.Background()))
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestRedis_Validate(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		defer c.Disconnect(context.Background())

		key := uuid.NewString()
		err = c.Validate(context.Background(), key)
		require.NoError(t, err)
	})

	t.Run("KO - key empty error", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		defer c.Disconnect(context.Background())

		err = c.Validate(context.Background(), "")
		require.ErrorIs(t, err, ErrKeyEmpty)
	})

	t.Run("KO - conflict error", func(st *testing.T) {
		c, err := NewRedis(redistestconf(), testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		defer c.Disconnect(context.Background())

		key := uuid.NewString()
		require.NoError(t, c.Validate(context.Background(), key))
		require.ErrorIs(t, c.Validate(context.Background(), key), ErrConflict)
	})
}

func redistestconf() *config.Config {
	testconf := &config.Config{
		Uri:        os.Getenv("REDIS_URI"),
		TimeToLive: testdata.Fake.UInt64Between(1000, 100000),
	}
	if testconf.Uri == "" {
		testconf.Uri = testdata.RedisUri
	}

	return testconf
}
