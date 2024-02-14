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

func (isntance *sql) Readiness() error {
	if isntance.status == patterns.StatusDisconnected {
		return nil
	}
	if isntance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	var ok int
	tx := isntance.client.Raw(ReadinessQuery).Scan(&ok)
	if tx.Error != nil || ok != 1 {
		return ErrNotReady
	}

	return nil
}

func (isntance *sql) Liveness() error {
	if isntance.status == patterns.StatusDisconnected {
		return nil
	}
	if isntance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	var ok int
	tx := isntance.client.Raw(LivenessQuery).Scan(&ok)
	if tx.Error != nil || ok != 1 {
		return ErrNotLive
	}

	return nil
}

func (isntance *sql) Connect(ctx context.Context) error {
	isntance.mu.Lock()
	defer isntance.mu.Unlock()

	if isntance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	client, err := NewGorm(isntance.conf, isntance.logger)
	if err != nil {
		return err
	}
	isntance.client = client

	isntance.status = patterns.StatusConnected
	isntance.logger.Info("connected")
	return nil
}

func (isntance *sql) Disconnect(ctx context.Context) error {
	isntance.mu.Lock()
	defer isntance.mu.Unlock()

	if isntance.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	isntance.status = patterns.StatusDisconnected
	isntance.logger.Info("disconnected")

	var returning error
	if conn, err := isntance.client.DB(); err == nil {
		if err := conn.Close(); err != nil {
			returning = errors.Join(returning, err)
		}
	} else {
		returning = errors.Join(returning, err)
	}
	isntance.client = nil

	return returning
}

func (isntance *sql) Client() any {
	return isntance.client
}
