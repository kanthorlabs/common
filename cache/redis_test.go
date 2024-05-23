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
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedis_New(t *testing.T) {
	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewRedis(conf)
		require.ErrorContains(t, err, "CACHE.CONFIG.")
	})
}

func TestRedis_Connect(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertConnect(t, cache, ErrAlreadyConnected)
}

func TestRedis_Readiness(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertReadiness(t, cache, ErrNotConnected)
}

func TestRedis_Liveness(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertLiveness(t, cache, ErrNotConnected)
}

func TestRedis_Disconnect(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertLiveness(t, cache, ErrNotConnected)
}

func TestRedis_Get(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, cache.Connect(ctx))
	defer cache.Disconnect(ctx)

	value := testdata.NewUser(clock.New())
	ttl := time.Minute

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		cache.Set(ctx, key, value, ttl)
		require.NoError(st, err)

		var dest testdata.User
		err := cache.Get(ctx, key, &dest)
		require.NoError(st, err)
		require.Equal(st, value, dest)
	})

	t.Run("KO - key of get method could not be empty", func(st *testing.T) {
		var dest string
		err := cache.Get(ctx, "", &dest)
		require.ErrorIs(st, err, ErrKeyEmpty)
	})

	t.Run("KO - key not found error", func(st *testing.T) {
		key := uuid.NewString()
		var dest testdata.User
		err := cache.Get(ctx, key, &dest)
		require.ErrorIs(st, err, ErrEntryNotFound)
	})

	t.Run("KO - unmarshal error", func(st *testing.T) {
		key := uuid.NewString()
		cache.Set(ctx, key, value, ttl)

		var dest chan int
		err := cache.Get(ctx, key, &dest)
		require.ErrorContains(st, err, "CACHE.VALUE.UNMARSHAL.ERROR")
	})
}

func TestRedis_Set(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)
	require.NoError(t, cache.Connect(ctx))
	defer cache.Disconnect(ctx)

	key := uuid.NewString()
	value := testdata.NewUser(clock.New())
	ttl := time.Minute

	t.Run("OK - not nil", func(st *testing.T) {
		key := uuid.NewString()
		err := cache.Set(ctx, key, value, ttl)
		require.NoError(st, err)
	})

	t.Run("OK - nil", func(st *testing.T) {
		err := cache.Set(ctx, key, nil, ttl)
		require.NoError(st, err)
	})

	t.Run("KO - key of set method could not be empty", func(st *testing.T) {
		err := cache.Set(ctx, "", value, ttl)
		require.ErrorIs(st, err, ErrKeyEmpty)
	})

	t.Run("KO - marshal error", func(st *testing.T) {
		err := cache.Set(ctx, key, make(chan int), ttl)
		require.ErrorContains(st, err, "CACHE.VALUE.MARSHAL.ERROR")
	})
}

func TestRedis_Exist(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)
	require.NoError(t, cache.Connect(ctx))
	defer cache.Disconnect(ctx)

	value := testdata.NewUser(clock.New())
	ttl := time.Minute

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()

		cache.Set(ctx, key, value, ttl)
		require.True(st, cache.Exist(ctx, key))
	})

	t.Run("KO - key of exist method could not be empty", func(st *testing.T) {
		require.False(st, cache.Exist(ctx, ""))
	})

	t.Run("OK - key not found err", func(st *testing.T) {
		key := uuid.NewString()

		require.False(st, cache.Exist(ctx, key))
	})
}

func TestRedis_Delete(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)
	require.NoError(t, cache.Connect(ctx))
	defer cache.Disconnect(ctx)

	value := testdata.NewUser(clock.New())
	ttl := time.Minute

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()

		cache.Set(ctx, key, value, ttl)
		require.True(st, cache.Exist(ctx, key))

		err := cache.Del(ctx, key)
		require.NoError(st, err)
		require.False(st, cache.Exist(ctx, key))
	})

	t.Run("KO - key of delete method could not be empty", func(st *testing.T) {
		err := cache.Del(ctx, "")
		require.ErrorIs(st, err, ErrKeyEmpty)
	})
}

func TestRedis_Expire(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)
	require.NoError(t, cache.Connect(ctx))
	defer cache.Disconnect(ctx)

	value := testdata.NewUser(clock.New())
	ttl := time.Second

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()

		cache.Set(ctx, key, value, ttl)
		require.True(st, cache.Exist(ctx, key))

		err := cache.Expire(ctx, key, time.Now().Add(time.Second))
		require.NoError(st, err)

		for i := 0; i < 10; i++ {
			if cache.Exist(ctx, key) {
				time.Sleep(time.Second / 2)
				continue
			}

			return
		}

		require.Fail(st, "key still exists after expiration")
	})

	t.Run("KO - key of expire method could not be empty", func(st *testing.T) {
		err := cache.Expire(ctx, "", time.Now())
		require.ErrorIs(st, err, ErrKeyEmpty)
	})

	t.Run("KO - key not found error", func(st *testing.T) {
		key := uuid.NewString()
		err := cache.Expire(ctx, key, time.Now().Add(time.Second))
		require.ErrorIs(st, err, ErrEntryNotFound)
	})

	t.Run("KO - negative ttl error", func(st *testing.T) {
		key := uuid.NewString()
		cache.Set(ctx, key, value, ttl)

		err := cache.Expire(ctx, key, time.Now().Add(-time.Second))
		require.ErrorContains(st, err, "CACHE.TIME_TO_LIVE.NEGATIVE.ERROR")
	})
}

func testConf(t *testing.T, container *redis.RedisContainer) *config.Config {
	uri, err := container.ConnectionString(context.Background())
	require.NoError(t, err)
	return &config.Config{Uri: uri}
}
