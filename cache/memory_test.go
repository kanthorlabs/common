package cache

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestMemory_New(t *testing.T) {
	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewMemory(conf, testify.Logger())
		require.ErrorContains(t, err, "CACHE.CONFIG.")
	})
}

func TestMemory_Connect(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)

	require.NoError(t, c.Connect(context.Background()))
	require.ErrorIs(t, c.Connect(context.Background()), ErrAlreadyConnected)
}

func TestMemory_Disconnect(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)

	require.ErrorIs(t, c.Disconnect(context.Background()), ErrNotConnected)
	require.NoError(t, c.Connect(context.Background()))
	require.NoError(t, c.Disconnect(context.Background()))
}

func TestMemory_Liveness(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)

	require.ErrorIs(t, c.Liveness(), ErrNotConnected)
	require.NoError(t, c.Connect(context.Background()))
	require.NoError(t, c.Liveness())
	require.NoError(t, c.Disconnect(context.Background()))
	require.NoError(t, c.Liveness())
}

func TestMemory_Readiness(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)

	require.ErrorIs(t, c.Readiness(), ErrNotConnected)
	require.NoError(t, c.Connect(context.Background()))
	require.NoError(t, c.Readiness())
	require.NoError(t, c.Disconnect(context.Background()))
	require.NoError(t, c.Readiness())
}

func TestMemory_Get(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)
	c.Connect(context.Background())
	defer c.Disconnect(context.Background())

	value := testdata.NewUser(clock.New())
	ttl := time.Minute

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		c.Set(context.Background(), key, value, ttl)

		var dest testdata.User
		err := c.Get(context.Background(), key, &dest)
		require.NoError(st, err)
		require.Equal(st, value, dest)
	})

	t.Run("KO - key is empty error", func(st *testing.T) {
		var dest string
		err := c.Get(context.Background(), "", &dest)
		require.ErrorContains(t, err, "CACHE.KEY.EMPTY.ERROR")
	})

	t.Run("KO - key not found error", func(st *testing.T) {
		key := uuid.NewString()
		var dest testdata.User
		err := c.Get(context.Background(), key, &dest)
		require.ErrorIs(st, err, ErrEntryNotFound)
	})

	t.Run("KO - unmarshal error", func(st *testing.T) {
		key := uuid.NewString()
		c.Set(context.Background(), key, value, ttl)

		var dest chan int
		err := c.Get(context.Background(), key, &dest)
		require.ErrorContains(st, err, "CACHE.VALUE.UNMARSHAL.ERROR")
	})
}

func TestMemory_Set(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)
	c.Connect(context.Background())
	defer c.Disconnect(context.Background())

	key := uuid.NewString()
	value := testdata.NewUser(clock.New())
	ttl := time.Minute

	t.Run("OK", func(st *testing.T) {
		err := c.Set(context.Background(), key, value, ttl)
		require.NoError(st, err)
	})

	t.Run("KO - key is empty error", func(st *testing.T) {
		err := c.Set(context.Background(), "", value, ttl)
		require.ErrorContains(st, err, "CACHE.KEY.EMPTY.ERROR")
	})

	t.Run("KO - marshal error", func(st *testing.T) {
		err := c.Set(context.Background(), key, make(chan int), ttl)
		require.ErrorContains(st, err, "CACHE.VALUE.MARSHAL.ERROR")
	})
}

func TestMemory_Exist(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)
	c.Connect(context.Background())
	defer c.Disconnect(context.Background())

	value := testdata.NewUser(clock.New())
	ttl := time.Minute

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()

		c.Set(context.Background(), key, value, ttl)
		require.True(st, c.Exist(context.Background(), key))
	})

	t.Run("KO - key is empty error", func(st *testing.T) {
		require.False(st, c.Exist(context.Background(), ""))
	})

	t.Run("OK - key not found err", func(st *testing.T) {
		key := uuid.NewString()

		require.False(st, c.Exist(context.Background(), key))
	})
}

func TestNenory_Del(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)
	c.Connect(context.Background())
	defer c.Disconnect(context.Background())

	value := testdata.NewUser(clock.New())
	ttl := time.Minute

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()

		c.Set(context.Background(), key, value, ttl)
		require.True(st, c.Exist(context.Background(), key))

		err := c.Del(context.Background(), key)
		require.NoError(st, err)
		require.False(st, c.Exist(context.Background(), key))
	})

	t.Run("KO - key is empty error", func(st *testing.T) {
		err := c.Del(context.Background(), "")
		require.ErrorContains(st, err, "CACHE.KEY.EMPTY.ERROR")
	})
}

func TestMemory_Expire(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)
	c.Connect(context.Background())
	defer c.Disconnect(context.Background())

	value := testdata.NewUser(clock.New())
	ttl := time.Second

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()

		c.Set(context.Background(), key, value, ttl)
		require.True(st, c.Exist(context.Background(), key))

		err := c.Expire(context.Background(), key, time.Now().Add(time.Second))
		require.NoError(st, err)

		for i := 0; i < 10; i++ {
			if c.Exist(context.Background(), key) {
				time.Sleep(time.Second / 2)
				continue
			}

			return
		}

		require.Fail(st, "key still exists after expiration")
	})

	t.Run("KO - key is empty error", func(st *testing.T) {
		err := c.Expire(context.Background(), "", time.Now())
		require.ErrorContains(st, err, "CACHE.KEY.EMPTY.ERROR")
	})

	t.Run("KO - key not found error", func(st *testing.T) {
		key := uuid.NewString()
		err := c.Expire(context.Background(), key, time.Now())
		require.ErrorIs(st, err, ErrEntryNotFound)
	})

	t.Run("KO - negative ttl error", func(st *testing.T) {
		key := uuid.NewString()
		c.Set(context.Background(), key, value, ttl)

		err := c.Expire(context.Background(), key, time.Now().Add(-time.Second))
		require.ErrorContains(st, err, "CACHE.TIME_TO_LIVE.NEGATIVE.ERROR")
	})
}
