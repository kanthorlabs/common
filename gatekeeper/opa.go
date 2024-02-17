package gatekeeper

import (
	"context"

	"github.com/kanthorlabs/common/gatekeeper/config"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/sqlx"
	"gorm.io/gorm"
)

func NewOpa(conf *config.Config, logger logging.Logger) (Gatekeeper, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	sequel, err := sqlx.New(&conf.Privilege.Sqlx, logger)
	if err != nil {
		return nil, err
	}

	return &opa{conf: conf, logger: logger, sequel: sequel}, nil
}

type opa struct {
	conf   *config.Config
	logger logging.Logger
	sequel persistence.Persistence

	orm *gorm.DB
}

func (instance *opa) Connect(ctx context.Context) error {
	if err := instance.sequel.Connect(ctx); err != nil {
		return err
	}

	instance.orm = instance.sequel.Client().(*gorm.DB)
	if err := instance.orm.WithContext(ctx).AutoMigrate(&entities.Privilege{}); err != nil {
		return err
	}

	return nil
}

func (instance *opa) Readiness() error {
	return instance.sequel.Readiness()
}

func (instance *opa) Liveness() error {
	return instance.sequel.Liveness()
}

func (instance *opa) Disconnect(ctx context.Context) error {
	return instance.sequel.Disconnect(ctx)
}

func (instance *opa) Grant(ctx context.Context, evaluation *entities.Evaluation) error {
	return nil
}
func (instance *opa) Enforce(ctx context.Context, evaluation *entities.Evaluation, permission *entities.Permission) error {
	return nil
}

func (instance *opa) Users(ctx context.Context, tenant string) ([]string, error) {
	return nil, nil
}

func (instance *opa) Tenants(ctx context.Context, username string) ([]string, error) {
	return nil, nil

}
