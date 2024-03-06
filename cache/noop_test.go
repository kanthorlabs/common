package cache

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNoop(t *testing.T) {
	cache := &noop{}
	ctx := context.Background()
	key := uuid.NewString()
	entry := testdata.NewUser(clock.New())
	ttl := time.Minute * time.Duration(testdata.Fake.Int64Between(1, 10))

	assert.NoError(t, cache.Connect(ctx))
	assert.NoError(t, cache.Readiness())
	assert.NoError(t, cache.Liveness())
	assert.ErrorIs(t, cache.Get(ctx, key, entry), ErrEntryNotFound)
	assert.NoError(t, cache.Set(ctx, key, entry, ttl))
	assert.False(t, cache.Exist(ctx, key))
	assert.NoError(t, cache.Del(ctx, key))
	assert.NoError(t, cache.Expire(ctx, key, time.Now()))
	assert.NoError(t, cache.Disconnect(ctx))
}
