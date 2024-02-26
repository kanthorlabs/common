package cache

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/mocks/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			conf := &config.Config{}
			_, err := NewMemory(conf, testify.Logger(), clock.NewClock(sst))
			require.ErrorContains(st, err, "CACHE.CONFIG.")
		})
	})

	t.Run(".Connect/.Readiness/.Liveness/.Disconnect", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger(), clock.NewClock(st))
		require.Nil(st, err)

		require.ErrorIs(st, c.Readiness(), ErrNotConnected)
		require.ErrorIs(st, c.Liveness(), ErrNotConnected)

		require.Nil(st, c.Connect(context.Background()))

		require.ErrorIs(st, c.Connect(context.Background()), ErrAlreadyConnected)

		require.Nil(st, c.Readiness())
		require.Nil(st, c.Liveness())

		require.Nil(st, c.Disconnect(context.Background()))

		require.Nil(st, c.Readiness())
		require.Nil(st, c.Liveness())

		require.ErrorIs(st, c.Disconnect(context.Background()), ErrNotConnected)
	})

	t.Run(".Get", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger(), clock.NewClock(st))
		require.Nil(st, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		st.Run("OK", func(sst *testing.T) {
			key := uuid.NewString()
			err := c.Set(context.Background(), key, testdata.Fake.Blood().Name(), time.Hour)
			require.Nil(st, err)

			_, err = c.Get(context.Background(), key)
			require.Nil(st, err)
		})

		st.Run("KO - not found error", func(sst *testing.T) {
			key := uuid.NewString()
			_, err := c.Get(context.Background(), key)
			require.ErrorIs(st, err, ErrEntryNotFound)
		})
	})

	t.Run(".Set", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger(), clock.NewClock(st))
		require.Nil(st, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		st.Run("OK - not nil", func(sst *testing.T) {
			key := uuid.NewString()
			value := testdata.Fake.Blood().Name()
			err := c.Set(context.Background(), key, value, time.Hour)
			require.Nil(st, err)
		})

		st.Run("OK - nil", func(sst *testing.T) {
			key := uuid.NewString()
			err := c.Set(context.Background(), key, nil, time.Hour)
			require.Nil(st, err)
		})

		st.Run("KO - marshal error", func(sst *testing.T) {
			key := uuid.NewString()
			err := c.Set(context.Background(), key, make(chan int), time.Hour)
			require.ErrorContains(st, err, "json: unsupported type")
		})
	})

	t.Run(".Exist", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger(), clock.NewClock(st))
		require.Nil(st, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		st.Run("OK", func(sst *testing.T) {
			key := uuid.NewString()
			exist := c.Exist(context.Background(), key)
			require.False(st, exist)
		})
	})

	t.Run(".Del", func(st *testing.T) {
		c, err := NewMemory(testconf, testify.Logger(), clock.NewClock(st))
		require.Nil(st, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		st.Run("OK", func(sst *testing.T) {
			key := uuid.NewString()
			err := c.Del(context.Background(), key)
			require.Nil(st, err)
		})
	})

	t.Run(".Expire", func(st *testing.T) {
		watch := clock.NewClock(st)
		c, err := NewMemory(testconf, testify.Logger(), watch)
		require.Nil(st, err)
		c.Connect(context.Background())
		defer c.Disconnect(context.Background())

		key := uuid.NewString()
		err = c.Set(context.Background(), key, testdata.Fake.Blood().Name(), time.Hour)
		require.Nil(st, err)

		st.Run("OK", func(sst *testing.T) {
			watch.EXPECT().Now().Return(time.Now().Add(-time.Hour)).Once()

			err = c.Expire(context.Background(), key, time.Now())
			require.Nil(st, err)
		})

		st.Run("KO - not found error", func(sst *testing.T) {
			err = c.Expire(context.Background(), uuid.NewString(), time.Now())
			require.ErrorIs(st, err, ErrEntryNotFound)
		})

		st.Run("KO - pass time", func(sst *testing.T) {
			watch.EXPECT().Now().Return(time.Now().Add(time.Hour)).Once()

			err = c.Expire(context.Background(), key, time.Now())
			require.ErrorContains(st, err, "CACHE.EXPIRE.EXPIRED_AT_TIME_PASS.ERROR")
		})
	})

}