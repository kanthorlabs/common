package sqlx

import (
	"context"
	"errors"
	"sync"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"gorm.io/gorm"
)

var (
	ReadinessQuery = "SELECT 1 as readiness"
	LivenessQuery  = "SELECT 1 as liveness"
)

func New(conf *config.Config, logger logging.Logger) (persistence.Persistence, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	logger = logger.With("storage", "sqlx")
	return &sql{conf: conf, logger: logger}, nil
}

type sql struct {
	conf   *config.Config
	logger logging.Logger

	client *gorm.DB

	mu     sync.Mutex
	status int
}

func (instance *sql) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	client, err := NewGorm(instance.conf, instance.logger)
	if err != nil {
		return err
	}
	instance.client = client

	instance.status = patterns.StatusConnected
	instance.logger.Info("connected")
	return nil
}

func (instance *sql) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	var ok int
	tx := instance.client.Raw(ReadinessQuery).Scan(&ok)
	if tx.Error != nil || ok != 1 {
		return ErrNotReady
	}

	return nil
}

func (instance *sql) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	var ok int
	tx := instance.client.Raw(LivenessQuery).Scan(&ok)
	if tx.Error != nil || ok != 1 {
		return ErrNotLive
	}

	return nil
}

func (instance *sql) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	instance.status = patterns.StatusDisconnected
	instance.logger.Info("disconnected")

	var returning error
	if conn, err := instance.client.DB(); err == nil {
		if err := conn.Close(); err != nil {
			returning = errors.Join(returning, err)
		}
	} else {
		returning = errors.Join(returning, err)
	}
	instance.client = nil

	return returning
}

func (instance *sql) Client() any {
	return instance.client
}
