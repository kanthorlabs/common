package distributedlockmanager

import (
	"context"
	"errors"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
)

func NewMemory(conf *config.Config) (Factory, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return func(key string, opts ...config.Option) DistributedLockManager {
		cconf := &config.Config{Uri: conf.Uri, TimeToLive: conf.TimeToLive}
		for _, opt := range opts {
			opt(cconf)
		}

		client := ttlcache.New[string, int]()
		go client.Start()
		return &memory{
			key:    key,
			client: client,
			conf:   cconf,
		}
	}, nil
}

type memory struct {
	key    string
	client *ttlcache.Cache[string, int]

	conf *config.Config
}

func (dlm *memory) Lock(ctx context.Context) error {
	k, err := Key(dlm.key)
	if err != nil {
		return err
	}

	if dlm.client.Has(k) {
		return errors.New("DISTRIBUTED_LOCK_MANAGER.LOCK.ERROR")
	}

	ttl := time.Millisecond * time.Duration(dlm.conf.TimeToLive)
	dlm.client.Set(k, int(1), ttl)
	return nil
}

func (dlm *memory) Unlock(ctx context.Context) error {
	k, err := Key(dlm.key)
	if err != nil {
		return err
	}

	if !dlm.client.Has(k) {
		return errors.New("DISTRIBUTED_LOCK_MANAGER.UNLOCK.ERROR")
	}

	dlm.client.Delete(k)
	return nil
}
