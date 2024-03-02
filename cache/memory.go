package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
)

// NewMemory creates a new memory cache instance that use ttlcache as the underlying engine.
func NewMemory(conf *config.Config, logger logging.Logger) (Cache, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	logger = logger.With("cache", "memory")
	return &memory{conf: conf, logger: logger}, nil
}

type memory struct {
	conf   *config.Config
	logger logging.Logger
	client *ttlcache.Cache[string, []byte]

	mu     sync.Mutex
	status int
}

func (instance *memory) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	instance.client = ttlcache.New[string, []byte]()
	go instance.client.Start()

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *memory) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return nil
}

func (instance *memory) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return nil
}

func (instance *memory) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	instance.status = patterns.StatusDisconnected

	instance.client.Stop()
	instance.client.DeleteAll()
	return nil
}

func (instance *memory) Get(ctx context.Context, key string, entry any) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	item := instance.client.Get(k)
	if item == nil {
		return ErrEntryNotFound
	}

	return Unmarshal(item.Value(), entry)
}

func (instance *memory) Set(ctx context.Context, key string, entry any, ttl time.Duration) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	v, err := Marshal(entry)
	if err != nil {
		return fmt.Errorf("CACHE.VALUE.MARSHAL.ERROR: %w", err)
	}

	instance.client.Set(k, v, ttl)
	return nil
}

func (instance *memory) Exist(ctx context.Context, key string) bool {
	k, err := Key(key)
	if err != nil {
		return false
	}

	return instance.client.Has(k)
}

func (instance *memory) Del(ctx context.Context, key string) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	instance.client.Delete(k)
	return nil
}

func (instance *memory) Expire(ctx context.Context, key string, at time.Time) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	item := instance.client.Get(k)
	if item == nil {
		return ErrEntryNotFound
	}

	ttl := time.Until(at)
	if ttl < 0 {
		return errors.New("CACHE.TIME_TO_LIVE.NEGATIVE.ERROR")
	}

	instance.client.Set(k, item.Value(), ttl)
	return nil
}
