package distributedlockmanager

import (
	"context"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/patterns"
)

// NewMemory creates a new distributed lock manager instance that uses ttlcache as the underlying storage.
func NewMemory(conf *config.Config) (DistributedLockManager, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	instance := &memory{conf: conf}
	return instance, nil
}

type memory struct {
	conf *config.Config

	client *ttlcache.Cache[string, int]

	mu     sync.Mutex
	status int
}

func (instance *memory) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	withTTL := ttlcache.WithTTL[string, int](time.Millisecond * time.Duration(instance.conf.TimeToLive))
	instance.client = ttlcache.New(withTTL)
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

func (instance *memory) Lock(ctx context.Context, key string, opts ...config.Option) (Identifier, error) {
	k, err := Key(key)
	if err != nil {
		return nil, err
	}

	conf := &config.Config{Uri: instance.conf.Uri, TimeToLive: instance.conf.TimeToLive}
	for _, opt := range opts {
		opt(conf)
	}

	if instance.client.Has(k) {
		return nil, ErrLock
	}

	ttl := time.Millisecond * time.Duration(instance.conf.TimeToLive)
	instance.client.Set(k, int(1), ttl)

	return &midentifier{k: k, client: instance.client}, nil
}

type midentifier struct {
	k      string
	client *ttlcache.Cache[string, int]
}

func (identifier *midentifier) Unlock(ctx context.Context) error {
	if !identifier.client.Has(identifier.k) {
		return ErrUnlock
	}

	identifier.client.Delete(identifier.k)
	return nil
}
