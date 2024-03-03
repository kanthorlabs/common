package distributedlockmanager

import (
	"context"
	"errors"
	"strings"

	"github.com/kanthorlabs/common/distributedlockmanager/config"
)

// New creates a new distributed lock manager instance based on the provided configuration.
// The distributed lock manager instance initialized based on the URI scheme.
// Supported schemes are:
// - memory://
// - redis://
// If the URI scheme is not supported, an error is returned.
func New(conf *config.Config) (Factory, error) {
	if strings.HasPrefix(conf.Uri, "memory") {
		return NewMemory(conf)
	}

	if strings.HasPrefix(conf.Uri, "redis") {
		return NewRedlock(conf)
	}

	return nil, errors.New("DISTRIBUTED_LOCK_MANAGER.SCHEME_UNKNOWN.ERROR")
}

type Factory func(key string, opts ...config.Option) DistributedLockManager

type DistributedLockManager interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}
