package strategies

import (
	"context"

	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/patterns"
)

type Strategy interface {
	patterns.Connectable
	Register(ctx context.Context, acc entities.Account) error
	Login(ctx context.Context, credentials entities.Credentials) (*entities.Tokens, error)
	Logout(ctx context.Context, tokens entities.Tokens) error
	Verify(ctx context.Context, tokens entities.Tokens) (*entities.Account, error)

	// Management APIs
	Deactivate(ctx context.Context, username string, at int64) error
	List(ctx context.Context, usernames []string) ([]*entities.Account, error)
	Update(ctx context.Context, account entities.Account) error
}
