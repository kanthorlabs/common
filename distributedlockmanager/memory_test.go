package distributedlockmanager

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestMemory_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := NewMemory(testconf)
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewMemory(conf)
		require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.CONFIG.")
	})
}

func TestMemory_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.ErrorIs(st, dlm.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestMemory_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Readiness())
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Disconnect(context.Background()))
		require.NoError(st, dlm.Readiness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.ErrorIs(st, dlm.Readiness(), ErrNotConnected)
	})
}

func TestMemory_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Liveness())
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Disconnect(context.Background()))
		require.NoError(st, dlm.Liveness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.ErrorIs(st, dlm.Liveness(), ErrNotConnected)
	})
}

func TestMemory_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.NoError(st, dlm.Connect(context.Background()))
		require.NoError(st, dlm.Disconnect(context.Background()))
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(t, err)

		require.ErrorIs(st, dlm.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestMemory_Lock(t *testing.T) {
	ttl := config.TimeToLive(testdata.Fake.UInt64Between(10000, 100000))
	dlm, err := NewMemory(testconf)
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

func TestMemory_Unlock(t *testing.T) {
	ttl := config.TimeToLive(testdata.Fake.UInt64Between(10000, 100000))
	dlm, err := NewMemory(testconf)
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
		k, _ := Key(uuid.NewString())
		identifier := &midentifier{k: k, client: dlm.(*memory).client}

		err = identifier.Unlock(context.Background())
		require.ErrorContains(st, err, ErrUnlock.Error())
	})
}
