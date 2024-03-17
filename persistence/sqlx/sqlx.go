package sqlx

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"gorm.io/gorm"
)

var (
	ReadinessQuery = "SELECT 1 as readiness"
	LivenessQuery  = "SELECT 1 as liveness"
)

func New(conf *config.Config, logger logging.Logger) (*SqlX, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	return &SqlX{conf: conf, logger: logger}, nil
}

type SqlX struct {
	conf   *config.Config
	logger logging.Logger

	client *gorm.DB

	mu     sync.Mutex
	status int
}

func (instance *SqlX) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	client, err := Gorm(instance.conf, instance.logger)
	if err != nil {
		return fmt.Errorf("SQLX.CONNECT.ERROR: %w", err)
	}
	instance.client = client

	instance.status = patterns.StatusConnected
	instance.logger.Info("connected")
	return nil
}

func (instance *SqlX) Readiness() error {
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

func (instance *SqlX) Liveness() error {
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

func (instance *SqlX) Disconnect(ctx context.Context) error {
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

func (instance *SqlX) Client() any {
	return instance.client
}
