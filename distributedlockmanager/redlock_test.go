package distributedlockmanager

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestRedlock_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := NewRedlock(redistestconf())
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewRedlock(conf)
		require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.CONFIG.")
	})
}

func TestRedlock_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.ErrorIs(st, dlm.Connect(context.Background()), ErrAlreadyConnected)
	})

	t.Run("KO - redis url error", func(st *testing.T) {
		conf := &config.Config{
			Uri:        "tcp://localhost:6379/0",
			TimeToLive: testdata.Fake.UInt64Between(10000, 100000),
		}
		dlm, err := NewRedlock(conf)
		require.NoError(t, err)

		require.ErrorContains(st, dlm.Connect(context.Background()), "redis: ")
	})
}

func TestRedlock_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Readiness())
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Disconnect(context.Background()))
		require.NoError(st, dlm.Readiness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.ErrorIs(st, dlm.Readiness(), ErrNotConnected)
	})
}

func TestRedlock_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Liveness())
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Disconnect(context.Background()))
		require.NoError(st, dlm.Liveness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.ErrorIs(st, dlm.Liveness(), ErrNotConnected)
	})
}

func TestRedlock_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Disconnect(context.Background()))
	})

	t.Run("KO", func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))

		// Close the connection to simulate the error of disconnection
		require.NoError(t, dlm.(*redlock).gredis.Close())

		require.ErrorContains(st, dlm.Disconnect(context.Background()), "redis: ")
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		dlm, err := NewRedlock(redistestconf())
		require.NoError(t, err)

		require.ErrorIs(st, dlm.Disconnect(context.Background()), ErrNotConnected)
	})

}

func TestRedlock_Lock(t *testing.T) {
	ttl := config.TimeToLive(testdata.Fake.UInt64Between(10000, 100000))
	dlm, err := NewRedlock(redistestconf())
	require.NoError(t, err)

	dlm.Connect(context.Background())
	defer dlm.Disconnect(context.Background())

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		identifier, err := dlm.Lock(context.Background(), key, ttl)
		require.NoError(st, err)
		require.NotNil(st, identifier)
	})

	t.Run(testify.CaseKOKeyEmptyError, func(st *testing.T) {
		_, err := dlm.Lock(context.Background(), "", ttl)
		require.ErrorIs(st, err, ErrKeyEmpty)
	})

	t.Run("KO - key already locked error", func(st *testing.T) {
		key := uuid.NewString()
		identifier, err := dlm.Lock(context.Background(), key, ttl)
		require.NoError(st, err)
		require.NotNil(st, identifier)

		_, err = dlm.Lock(context.Background(), key, ttl)
		require.ErrorContains(st, err, ErrLock.Error())
	})
}

func TestRedlock_Unlock(t *testing.T) {
	ttl := config.TimeToLive(testdata.Fake.UInt64Between(10000, 100000))
	dlm, err := NewRedlock(redistestconf())
	require.NoError(t, err)

	dlm.Connect(context.Background())
	defer dlm.Disconnect(context.Background())

	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		identifier, err := dlm.Lock(context.Background(), key, ttl)
		require.NoError(st, err)
		require.NotNil(st, identifier)

		require.NoError(st, identifier.Unlock(context.Background()))
	})

	t.Run("KO - key not locked error", func(st *testing.T) {
		key := uuid.NewString()
		identifier, err := dlm.Lock(context.Background(), key, ttl)
		require.NoError(st, err)
		require.NotNil(st, identifier)

		require.NoError(st, identifier.Unlock(context.Background()))
		require.ErrorContains(st, identifier.Unlock(context.Background()), ErrUnlock.Error())

	})
}

func redistestconf() *config.Config {
	testconf := &config.Config{
		Uri:        os.Getenv("REDIS_URI"),
		TimeToLive: testdata.Fake.UInt64Between(10000, 100000),
	}
	if testconf.Uri == "" {
		testconf.Uri = testdata.RedisUri
	}

	return testconf
}
