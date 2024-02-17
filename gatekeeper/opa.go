package gatekeeper

import (
	"context"
	"errors"
	"time"

	"github.com/kanthorlabs/common/gatekeeper/config"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/sqlx"
	"github.com/kanthorlabs/common/safe"
	"gorm.io/gorm"
)

func New(conf *config.Config, logger logging.Logger) (Gatekeeper, error) {
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
	privilege := &entities.Privilege{
		Tenant:    evaluation.Tenant,
		Username:  evaluation.Username,
		Role:      evaluation.Role,
		Metadata:  &safe.Metadata{},
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	tx := instance.orm.WithContext(ctx).Create(privilege)

	return tx.Error
}

func (instance *opa) Revoke(ctx context.Context, evaluation *entities.Evaluation) error {
	privilege := &entities.Privilege{
		Tenant:   evaluation.Tenant,
		Username: evaluation.Username,
		Role:     evaluation.Role,
	}
	tx := instance.orm.WithContext(ctx).
		Where(
			"tenant = ? AND username = ? AND role = ?",
			evaluation.Tenant, evaluation.Username, evaluation.Role,
		).
		Delete(privilege)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("GATEKEEPER.REVOKE.PRIVILEGE_NOT_EXIST.ERROR")
	}

	return nil
}

func (instance *opa) Enforce(ctx context.Context, evaluation *entities.Evaluation, permission *entities.Permission) error {
	return nil
}

func (instance *opa) Users(ctx context.Context, tenant string) ([]entities.User, error) {
	var privileges []entities.Privilege

	tx := instance.orm.WithContext(ctx).
		Model(&entities.Privilege{}).
		Where("tenant = ?", tenant).
		Find(&privileges)

	if tx.Error != nil {
		return nil, tx.Error
	}

	maps := map[string][]string{}
	for _, privilege := range privileges {
		if _, exist := maps[privilege.Username]; !exist {
			maps[privilege.Username] = []string{privilege.Role}
			continue
		}

		maps[privilege.Username] = append(maps[privilege.Username], privilege.Role)
	}

	users := make([]entities.User, 0)
	for username := range maps {
		users = append(users, entities.User{
			Username: username,
			Roles:    maps[username],
		})
	}

	return users, nil
}

func (instance *opa) Tenants(ctx context.Context, username string) ([]entities.Tenant, error) {
	var privileges []entities.Privilege
	tx := instance.orm.WithContext(ctx).
		Model(&entities.Privilege{}).
		Where("username = ?", username).
		Find(&privileges)

	if tx.Error != nil {
		return nil, tx.Error
	}

	maps := map[string][]string{}
	for _, privilege := range privileges {
		if _, exist := maps[privilege.Tenant]; !exist {
			maps[privilege.Tenant] = []string{privilege.Role}
			continue
		}

		maps[privilege.Tenant] = append(maps[privilege.Tenant], privilege.Role)
	}

	tenants := make([]entities.Tenant, 0)
	for tenant := range maps {
		tenants = append(tenants, entities.Tenant{
			Tenant: tenant,
			Roles:  maps[tenant],
		})
	}

	return tenants, nil
}
