package distributedlockmanager

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedlock_New(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		_, err := NewRedlock(testConf(t, container))
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := NewRedlock(&config.Config{})
		require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.CONFIG.")
	})
}

func TestRedlock_Connect(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)

	dlm, err := NewRedlock(testConf(t, container))
	require.NoError(t, err)

	testify.AssertConnect(t, dlm, ErrAlreadyConnected)
}

func TestRedlock_Readiness(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)

	dlm, err := NewRedlock(testConf(t, container))
	require.NoError(t, err)

	testify.AssertReadiness(t, dlm, ErrNotConnected)
}

func TestRedlock_Liveness(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)

	dlm, err := NewRedlock(testConf(t, container))
	require.NoError(t, err)

	testify.AssertLiveness(t, dlm, ErrNotConnected)
}

func TestRedlock_Disconnect(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)

	dlm, err := NewRedlock(testConf(t, container))
	require.NoError(t, err)

	testify.AssertLiveness(t, dlm, ErrNotConnected)
}

func TestRedlock_Lock(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)

	dlm, err := NewRedlock(testConf(t, container))
	require.NoError(t, err)

	require.NoError(t, dlm.Connect(ctx))
	defer dlm.Disconnect(context.Background())

	ttl := config.TimeToLive(testdata.Fake.UInt64Between(10000, 100000))

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		identifier, err := dlm.Lock(context.Background(), key, ttl)
		require.NoError(st, err)
		require.NotNil(st, identifier)
	})

	t.Run("KO - lock key empty error", func(st *testing.T) {
		_, err := dlm.Lock(context.Background(), "", ttl)
		require.ErrorIs(st, err, ErrKeyEmpty)
	})

	t.Run("KO - lock key already locked error", func(st *testing.T) {
		key := uuid.NewString()
		identifier, err := dlm.Lock(context.Background(), key, ttl)
		require.NoError(st, err)
		require.NotNil(st, identifier)

		_, err = dlm.Lock(context.Background(), key, ttl)
		require.ErrorContains(st, err, ErrLock.Error())
	})
}

func TestRedlock_Unlock(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)

	dlm, err := NewRedlock(testConf(t, container))
	require.NoError(t, err)

	require.NoError(t, dlm.Connect(ctx))
	defer dlm.Disconnect(context.Background())

	ttl := config.TimeToLive(testdata.Fake.UInt64Between(10000, 100000))

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		identifier, err := dlm.Lock(context.Background(), key, ttl)
		require.NoError(st, err)
		require.NotNil(st, identifier)

		require.NoError(st, identifier.Unlock(context.Background()))
	})

	t.Run("KO - lock key is not locked error", func(st *testing.T) {
		key := uuid.NewString()
		identifier, err := dlm.Lock(context.Background(), key, ttl)
		require.NoError(st, err)
		require.NotNil(st, identifier)

		require.NoError(st, identifier.Unlock(context.Background()))
		require.ErrorContains(st, identifier.Unlock(context.Background()), ErrUnlock.Error())
	})
}

func testConf(t *testing.T, container *redis.RedisContainer) *config.Config {
	uri, err := container.ConnectionString(context.Background())
	require.NoError(t, err)
	return &config.Config{Uri: uri, TimeToLive: testdata.Fake.UInt64Between(10000, 100000)}
}
