package strategies

import (
	"context"

	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/patterns"
)

type Strategy interface {
	patterns.Connectable
	Login(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error)
	Logout(ctx context.Context, credentials *entities.Credentials) error
	Verify(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error)
	Register(ctx context.Context, acc *entities.Account) (*entities.Account, error)
}
