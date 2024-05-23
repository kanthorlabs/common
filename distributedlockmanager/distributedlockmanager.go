package distributedlockmanager

import (
	"context"
	"errors"
	"strings"

	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/patterns"
)

func New(conf *config.Config) (DistributedLockManager, error) {
	if strings.HasPrefix(conf.Uri, "redis") {
		return NewRedlock(conf)
	}

	return nil, errors.New("DISTRIBUTED_LOCK_MANAGER.SCHEME_UNKNOWN.ERROR")
}

type DistributedLockManager interface {
	patterns.Connectable
	Lock(ctx context.Context, key string, opts ...config.Option) (Identifier, error)
}

type Identifier interface {
	Unlock(ctx context.Context) error
}
