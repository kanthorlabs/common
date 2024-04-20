package idempotency

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestMemory_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := NewMemory(conf, testify.Logger())
		require.ErrorContains(t, err, "IDEMPOTENCY.CONFIG.")
	})
}

func TestMemory_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		require.ErrorIs(t, c.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestMemory_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Readiness(), ErrNotConnected)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Readiness())
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Disconnect(context.Background()))

		require.NoError(t, c.Readiness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Readiness(), ErrNotConnected)
	})
}

func TestMemory_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Liveness(), ErrNotConnected)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Liveness())
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Disconnect(context.Background()))

		require.NoError(t, c.Liveness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Liveness(), ErrNotConnected)
	})
}

func TestMemory_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, c.Connect(context.Background()))
		require.NoError(t, c.Disconnect(context.Background()))
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, c.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestMemory_Validate(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		key := uuid.NewString()
		require.NoError(t, c.Validate(context.Background(), key))
	})

	t.Run("KO - key empty error", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		require.ErrorIs(t, c.Validate(context.Background(), ""), ErrKeyEmpty)
	})

	t.Run("KO - conflict error", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(t, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		key := uuid.NewString()
		require.NoError(t, c.Validate(context.Background(), key))
		require.ErrorIs(t, c.Validate(context.Background(), key), ErrConflict)
	})
}
