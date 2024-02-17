package gatekeeper

import (
	"context"

	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/patterns"
)

// Gatekeeper is an implementation of multi-tenant RBAC
type Gatekeeper interface {
	patterns.Connectable
	Grant(ctx context.Context, evaluation *entities.Evaluation) error
	Enforce(ctx context.Context, evaluation *entities.Evaluation, permission *entities.Permission) error
	Users(ctx context.Context, tenant string) ([]string, error)
	Tenants(ctx context.Context, username string) ([]string, error)
}
