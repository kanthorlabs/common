package gatekeeper

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kanthorlabs/common/gatekeeper/config"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/gatekeeper/rego"
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

	mu          sync.Mutex
	orm         *gorm.DB
	definitions map[string][]entities.Permission
	evaluate    rego.Evaluate
}

func (instance *opa) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if err := instance.sequel.Connect(ctx); err != nil {
		return err
	}

	instance.orm = instance.sequel.Client().(*gorm.DB)
	if err := instance.orm.WithContext(ctx).AutoMigrate(&entities.Privilege{}); err != nil {
		return err
	}

	definitions, err := config.ParseDefinitionsToPermissions(instance.conf.Definitions.Uri)
	if err != nil {
		return err
	}
	instance.definitions = definitions

	evaluate, err := rego.RBAC(ctx, definitions)
	if err != nil {
		return err
	}
	instance.evaluate = evaluate

	return nil
}

func (instance *opa) Readiness() error {
	return instance.sequel.Readiness()
}

func (instance *opa) Liveness() error {
	return instance.sequel.Liveness()
}

func (instance *opa) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	return instance.sequel.Disconnect(ctx)
}

func (instance *opa) Grant(ctx context.Context, evaluation *entities.Evaluation) error {
	if err := entities.EvaluationValidateOnGrant(evaluation); err != nil {
		return err
	}

	if _, exist := instance.definitions[evaluation.Role]; !exist {
		return errors.New("GATEKEEPER.GRANT.ROLE_NOT_EXIST.ERROR")
	}

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
	if err := entities.EvaluationValidateOnRevoke(evaluation); err != nil {
		return err
	}

	tx := instance.orm.WithContext(ctx).
		Where(
			"tenant = ? AND username = ?",
			evaluation.Tenant, evaluation.Username,
		).
		Delete(&entities.Privilege{})
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("GATEKEEPER.REVOKE.PRIVILEGE_NOT_EXIST.ERROR")
	}

	return nil
}

// Enforce will attempt to evaluate all roles, so there's no need to use evaluation.Role for validation or querying.
func (instance *opa) Enforce(ctx context.Context, evaluation *entities.Evaluation, permission *entities.Permission) error {
	if err := entities.EvaluationValidateOnEnforce(evaluation); err != nil {
		return err
	}
	if err := permission.Validate(); err != nil {
		return err
	}

	var privileges []entities.Privilege
	tx := instance.orm.WithContext(ctx).
		Model(&entities.Privilege{}).
		Where("tenant = ? AND username = ?", evaluation.Tenant, evaluation.Username).
		Find(&privileges)

	if tx.Error != nil {
		return tx.Error
	}

	if len(privileges) == 0 {
		return errors.New("GATEKEEPER.ENFORCE.PRIVILEGE_EMPTY.ERROR")
	}

	return instance.evaluate(permission, privileges)
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
