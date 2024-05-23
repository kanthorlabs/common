package cache

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"strings"
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

func TestEncodeKey(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		count := testdata.Fake.IntBetween(10, 100)
		keys := make([]string, count)
		for i := 0; i < count; i++ {
			keys[i] = uuid.NewString()
		}

		h := sha256.New()
		io.WriteString(h, strings.Join(keys, "/"))
		expected := fmt.Sprintf("%x", h.Sum(nil))

		require.Equal(st, expected, EncodeKey(keys...))
	})

	t.Run("OK - empty keys should return value", func(st *testing.T) {
		keys := make([]string, 0)
		h := sha256.New()
		io.WriteString(h, strings.Join(keys, "/"))
		expected := fmt.Sprintf("%x", h.Sum(nil))

		require.Equal(st, expected, EncodeKey(keys...))
	})
}

func TestGetOrSet(t *testing.T) {
	ctx := context.Background()
	container, err := testify.RedisContainer(ctx)
	require.NoError(t, err)

	cache, err := NewRedis(testConf(t, container))
	require.NoError(t, err)
	require.NoError(t, cache.Connect(ctx))
	defer cache.Disconnect(ctx)

	ttl := time.Minute
	value := testdata.NewUser(clock.New())

	t.Run("OK - get from cache", func(st *testing.T) {
		key := uuid.NewString()
		err := cache.Set(context.Background(), key, value, ttl)
		require.NoError(st, err)

		entry, err := GetOrSet(cache, context.Background(), key, ttl, func() (*testdata.User, error) {
			return nil, nil
		})

		require.NoError(st, err)
		require.Equal(st, value, *entry)
	})

	t.Run("OK - get from fn", func(st *testing.T) {
		key := uuid.NewString()

		entry, err := GetOrSet(cache, context.Background(), key, ttl, func() (*testdata.User, error) {
			return &value, nil
		})

		require.NoError(st, err)
		require.Equal(st, value, *entry)
	})

	t.Run("KO - get from cache error", func(st *testing.T) {
		key := uuid.NewString()
		err := cache.Set(context.Background(), key, "-", ttl)
		require.NoError(st, err)

		_, err = GetOrSet(cache, context.Background(), key, 0, func() (*testdata.User, error) {
			return nil, nil
		})

		require.Error(st, err)
	})

	t.Run("KO - get from fn error", func(st *testing.T) {
		key := uuid.NewString()
		expected := errors.New("error")

		_, err := GetOrSet(cache, context.Background(), key, 0, func() (*string, error) {
			return nil, expected
		})

		require.ErrorIs(st, err, expected)
	})
}
