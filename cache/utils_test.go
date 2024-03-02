package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		key := uuid.NewString()
		expected := "cache/" + key

		k, err := Key(key)

		require.NoError(st, err)
		require.Equal(st, expected, k)
	})

	t.Run("KO - empty key error", func(st *testing.T) {
		k, err := Key("")

		require.Empty(st, k)
		require.ErrorIs(st, err, ErrKeyEmpty)
	})
}

func TestGetOrSet(t *testing.T) {
	c, err := NewMemory(testconf, testify.Logger())
	require.NoError(t, err)
	c.Connect(context.Background())
	defer c.Disconnect(context.Background())

	ttl := time.Minute
	value := testdata.NewUser(clock.New())

	t.Run("OK - get from cache", func(st *testing.T) {
		key := uuid.NewString()
		err := c.Set(context.Background(), key, value, ttl)
		require.NoError(st, err)

		entry, err := GetOrSet(c, context.Background(), key, ttl, func() (*testdata.User, error) {
			return nil, nil
		})

		require.NoError(st, err)
		require.Equal(st, value, *entry)
	})

	t.Run("OK - get from fn", func(st *testing.T) {
		key := uuid.NewString()

		entry, err := GetOrSet(c, context.Background(), key, ttl, func() (*testdata.User, error) {
			return &value, nil
		})

		require.NoError(st, err)
		require.Equal(st, value, *entry)
	})

	t.Run("KO - get from cache error", func(st *testing.T) {
		key := uuid.NewString()
		err := c.Set(context.Background(), key, "-", ttl)
		require.NoError(st, err)

		_, err = GetOrSet(c, context.Background(), key, 0, func() (*testdata.User, error) {
			return nil, nil
		})

		require.Error(st, err)
	})

	t.Run("KO - get from fn error", func(st *testing.T) {
		key := uuid.NewString()
		expected := errors.New("error")

		_, err := GetOrSet(c, context.Background(), key, 0, func() (*string, error) {
			return nil, expected
		})

		require.ErrorIs(st, err, expected)
	})
}
