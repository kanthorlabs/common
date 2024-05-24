package idempotency

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/containers"
	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedis_New(t *testing.T) {
	ctx := context.Background()
	container, err := containers.Redis(ctx, "kanthorlabs-common-idempotency")
	require.NoError(t, err)

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
	container, err := containers.Redis(ctx, "kanthorlabs-common-idempotency")
	require.NoError(t, err)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertConnect(t, idemp, ErrAlreadyConnected)
}

func TestRedis_Readiness(t *testing.T) {
	ctx := context.Background()
	container, err := containers.Redis(ctx, "kanthorlabs-common-idempotency")
	require.NoError(t, err)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertReadiness(t, idemp, ErrNotConnected)
}

func TestRedis_Liveness(t *testing.T) {
	ctx := context.Background()
	container, err := containers.Redis(ctx, "kanthorlabs-common-idempotency")
	require.NoError(t, err)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertLiveness(t, idemp, ErrNotConnected)
}

func TestRedis_Disconnect(t *testing.T) {
	ctx := context.Background()
	container, err := containers.Redis(ctx, "kanthorlabs-common-idempotency")
	require.NoError(t, err)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, err)

	testify.AssertLiveness(t, idemp, ErrNotConnected)
}

func TestRedis_Validate(t *testing.T) {
	ctx := context.Background()
	container, err := containers.Redis(ctx, "kanthorlabs-common-idempotency")
	require.NoError(t, err)

	idemp, err := NewRedis(testConf(t, container))
	require.NoError(t, idemp.Connect(ctx))
	defer idemp.Disconnect(ctx)

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		require.NoError(st, idemp.Validate(ctx, key))
	})

	t.Run("KO - key of validate method could not be empty", func(st *testing.T) {
		err = idemp.Validate(context.Background(), "")
		require.ErrorIs(st, err, ErrKeyEmpty)
	})

	t.Run("KO - conflict error", func(st *testing.T) {
		key := uuid.NewString()
		require.NoError(st, idemp.Validate(context.Background(), key))
		require.ErrorIs(st, idemp.Validate(context.Background(), key), ErrConflict)
	})
}

func testConf(t *testing.T, container *redis.RedisContainer) *config.Config {
	uri, err := containers.RedisConnectionString(context.Background(), container)
	require.NoError(t, err)
	return &config.Config{Uri: uri, TimeToLive: testdata.Fake.UInt64Between(10000, 100000)}
}
