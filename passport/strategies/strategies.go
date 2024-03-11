package strategies

import (
	"context"

	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/patterns"
)

type Strategy interface {
	patterns.Connectable
	ParseCredentials(ctx context.Context, raw string) (*entities.Credentials, error)
	Login(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error)
	Logout(ctx context.Context, credentials *entities.Credentials) error
	Verify(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error)
	Register(ctx context.Context, acc *entities.Account) error
	Deactivate(ctx context.Context, username string, at int64) error

	// Management APIs
	List(ctx context.Context, usernames []string) ([]*entities.Account, error)
	Update(ctx context.Context, account *entities.Account) error
}
