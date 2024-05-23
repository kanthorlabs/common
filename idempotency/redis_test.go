package idempotency

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedis_New(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	t.Run("OK", func(st *testing.T) {
		_, err := NewRedis(testConf(t, container))
		require.NoError(t, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewRedis(conf)
		require.ErrorContains(t, err, "IDEMPOTENCY.CONFIG.")
	})
}

func TestRedis_Connect(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertConnect(t, idemp, ErrAlreadyConnected)
}

func TestRedis_Readiness(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertReadiness(t, idemp, ErrNotConnected)
}

func TestRedis_Liveness(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertLiveness(t, idemp, ErrNotConnected)
}

func TestRedis_Disconnect(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertLiveness(t, idemp, ErrNotConnected)
}

func TestRedis_Validate(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, idemp.Connect(ctx))
	defer idemp.Disconnect(ctx)

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		require.NoError(t, idemp.Validate(ctx, key))
	})

	t.Run("KO - key of validate method could not be empty", func(st *testing.T) {
		err = idemp.Validate(context.Background(), "")
		require.ErrorIs(t, err, ErrKeyEmpty)
	})

	t.Run("KO - conflict error", func(st *testing.T) {
		key := uuid.NewString()
		require.NoError(t, idemp.Validate(context.Background(), key))
		require.ErrorIs(t, idemp.Validate(context.Background(), key), ErrConflict)
	})
}

func testConf(t *testing.T, container *redis.RedisContainer) *config.Config {
	uri, err := container.ConnectionString(context.Background())
	require.NoError(t, err)
	return &config.Config{Uri: uri, TimeToLive: testdata.Fake.UInt64Between(10000, 100000)}
}
