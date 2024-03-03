package idempotency

import (
	"context"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
)

func NewMemory(conf *config.Config, logger logging.Logger) (Idempotency, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	logger = logger.With("idempotency", "memory")
	return &memory{conf: conf, logger: logger}, nil
}

type memory struct {
	conf   *config.Config
	logger logging.Logger
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

	instance.client = ttlcache.New[string, int]()
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

func (instance *memory) Validate(ctx context.Context, key string) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	if instance.client.Has(k) {
		return ErrConflict
	}

	ttl := time.Millisecond * time.Duration(instance.conf.TimeToLive)
	instance.client.Set(k, int(1), ttl)
	return nil
}
