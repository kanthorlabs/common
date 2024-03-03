package gatekeeper

import (
	"context"

	"github.com/kanthorlabs/common/gatekeeper/config"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
)

func New(conf *config.Config, logger logging.Logger) (Gatekeeper, error) {
	return NewOpa(conf, logger)
}

// Gatekeeper is an implementation of multi-tenant RBAC
type Gatekeeper interface {
	patterns.Connectable
	Grant(ctx context.Context, evaluation *entities.Evaluation) error
	Revoke(ctx context.Context, evaluation *entities.Evaluation) error
	Enforce(ctx context.Context, evaluation *entities.Evaluation, permission *entities.Permission) error
	Users(ctx context.Context, tenant string) ([]entities.User, error)
	Tenants(ctx context.Context, username string) ([]entities.Tenant, error)
}
