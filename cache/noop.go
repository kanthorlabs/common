package cache

import (
	"context"
	"time"
)

func NewNoop() Cache {
	return &noop{}
}

type noop struct{}

func (cache *noop) Connect(ctx context.Context) error {
	return nil
}

func (cache *noop) Readiness() error {
	return nil
}

func (cache *noop) Liveness() error {
	return nil
}

func (cache *noop) Disconnect(ctx context.Context) error {
	return nil
}

func (cache *noop) Get(ctx context.Context, key string, entry any) error {
	return ErrEntryNotFound
}

func (cache *noop) Set(ctx context.Context, key string, entry any, ttl time.Duration) error {
	return nil
}

func (cache *noop) Exist(ctx context.Context, key string) bool {
	return false
}

func (cache *noop) Del(ctx context.Context, key string) error {
	return nil
}

func (cache *noop) Expire(ctx context.Context, key string, at time.Time) error {
	return nil
}
