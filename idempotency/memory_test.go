package idempotency

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			conf := &config.Config{}
			_, err := NewMemory(conf, testify.Logger())
			require.ErrorContains(st, err, "IDEMPOTENCY.CONFIG.")
		})
	})

	t.Run(".Connect/.Readiness/.Liveness/.Disconnect", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Readiness(), ErrNotConnected)
		require.ErrorIs(st, c.Liveness(), ErrNotConnected)

		require.NoError(st, c.Connect(context.Background()))

		require.ErrorIs(st, c.Connect(context.Background()), ErrAlreadyConnected)

		require.NoError(st, c.Readiness())
		require.NoError(st, c.Liveness())

		require.NoError(st, c.Disconnect(context.Background()))

		require.NoError(st, c.Readiness())
		require.NoError(st, c.Liveness())

		require.ErrorIs(st, c.Disconnect(context.Background()), ErrNotConnected)
	})

	t.Run(".Validate", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger())
		require.NoError(st, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		st.Run("OK", func(sst *testing.T) {
			key := uuid.NewString()
			err := c.Validate(context.Background(), key)
			require.NoError(st, err)
		})

		st.Run("KO", func(sst *testing.T) {
			key := uuid.NewString()
			require.NoError(st, c.Validate(context.Background(), key))
			require.ErrorIs(st, c.Validate(context.Background(), key), ErrConflict)
		})
	})
}
