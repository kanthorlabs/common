package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
)

func NewMemory(conf *config.Config, logger logging.Logger, watch clock.Clock) (Cache, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	cache := ttlcache.New[string, []byte]()
	return &memory{cache: cache, conf: conf, logger: logger, watch: watch}, nil
}

type memory struct {
	conf   *config.Config
	logger logging.Logger
	watch  clock.Clock
	cache  *ttlcache.Cache[string, []byte]

	mu     sync.Mutex
	status int
}

func (instance *memory) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	go instance.cache.Start()

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

	instance.cache.Stop()
	instance.cache.DeleteAll()
	return nil
}

func (instance *memory) Get(ctx context.Context, key string) ([]byte, error) {
	item := instance.cache.Get(key)
	if item == nil {
		return nil, ErrEntryNotFound
	}
	return item.Value(), nil
}

func (instance *memory) Set(ctx context.Context, key string, entry any, ttl time.Duration) error {
	var value []byte
	var err error
	if entry != nil {
		value, err = json.Marshal(entry)
		if err != nil {
			return err
		}
	}
	instance.cache.Set(key, value, ttl)
	return nil
}

func (instance *memory) Exist(ctx context.Context, key string) bool {
	return instance.cache.Has(key)
}

func (instance *memory) Del(ctx context.Context, key string) error {
	instance.cache.Delete(key)
	return nil
}

func (instance *memory) Expire(ctx context.Context, key string, at time.Time) error {
	value, err := instance.Get(ctx, key)
	if err != nil {
		return err
	}

	ttl := at.Sub(instance.watch.Now())
	if ttl < 0 {
		return errors.New("CACHE.EXPIRE.EXPIRED_AT_TIME_PASS.ERROR")
	}

	instance.cache.Set(key, value, ttl)
	return nil
}
