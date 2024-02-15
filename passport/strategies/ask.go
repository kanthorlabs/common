package strategies

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
)

func NewAsk(conf *config.Ask, logger logging.Logger) (Strategy, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	accounts := make(map[string]*entities.Account)
	for i := range conf.Accounts {
		accounts[conf.Accounts[i].Username] = &conf.Accounts[i]
	}

	return &ask{conf: conf, logger: logger, accounts: accounts}, nil
}

type ask struct {
	conf   *config.Ask
	logger logging.Logger

	accounts map[string]*entities.Account
}

func (instance *ask) Readiness() error {
	return nil
}

func (instance *ask) Liveness() error {
	return nil
}

func (instance *ask) Connect(ctx context.Context) error {
	return nil
}

func (instance *ask) Disconnect(ctx context.Context) error {
	return nil
}

func (instance *ask) Login(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error) {
	if err := entities.ValidateCredentialsOnLogin(credentials); err != nil {
		return nil, err
	}
	acc, has := instance.accounts[credentials.Username]
	if !has {
		return nil, ErrLogin
	}

	if err := password.CompareString(credentials.Password, acc.PasswordHash); err != nil {
		return nil, ErrLogin
	}

	return acc.Censor(), nil
}

func (instance *ask) Logout(ctx context.Context, credentials *entities.Credentials) error {
	return nil
}

func (instance *ask) Verify(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error) {
	return instance.Login(ctx, credentials)
}

func (instance *ask) Register(ctx context.Context, acc *entities.Account) error {
	return errors.New("PASSPORT.ASK.REGISTER.UNIMPLEMENT.ERROR")
}
