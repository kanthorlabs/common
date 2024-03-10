package strategies

import (
	"context"
	"errors"
	"sync"

	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/patterns"
)

// NewAsk creates a new ask strategy instance what allows to authenticate users based on a username and password.
// The storage is based on a map of accounts in memory.
func NewAsk(conf *config.Ask, logger logging.Logger) (Strategy, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	accounts := make(map[string]*entities.Account)
	for i := range conf.Accounts {
		accounts[conf.Accounts[i].Username] = &conf.Accounts[i]
	}

	if len(accounts) != len(conf.Accounts) {
		return nil, errors.New("PASSPORT.STRATEGY.ASK.DUPLICATED_ACCOUNT.ERROR")
	}

	return &ask{conf: conf, logger: logger, accounts: accounts}, nil
}

type ask struct {
	conf   *config.Ask
	logger logging.Logger

	mu       sync.Mutex
	status   int
	accounts map[string]*entities.Account
}

func (instance *ask) ParseCredentials(ctx context.Context, raw string) (*entities.Credentials, error) {
	if IsBasicScheme(raw) {
		return ParseBasicCredentials(raw)
	}

	return nil, ErrCredentialsScheme
}

func (instance *ask) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *ask) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return nil
}

func (instance *ask) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return nil
}

func (instance *ask) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	instance.status = patterns.StatusDisconnected
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

	if err := password.Compare(credentials.Password, acc.PasswordHash); err != nil {
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

func (instance *ask) Deactivate(ctx context.Context, username string, ts int64) error {
	return errors.New("PASSPORT.ASK.DEACTIVATE.UNIMPLEMENT.ERROR")
}
