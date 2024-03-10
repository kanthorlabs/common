package strategies

import (
	"context"
	"sync"
	"time"

	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/sqlx"
	"gorm.io/gorm"
)

// NewDurability creates a new durability strategy instance what allows to authenticate users based on a username and password.
// The storage is based on a SQL database.
func NewDurability(conf *config.Durability, logger logging.Logger) (Strategy, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	sequel, err := sqlx.New(&conf.Sqlx, logger)
	if err != nil {
		return nil, err
	}

	return &durability{conf: conf, logger: logger, sequel: sequel}, nil
}

type durability struct {
	conf   *config.Durability
	logger logging.Logger
	sequel persistence.Persistence

	mu     sync.Mutex
	status int
	orm    *gorm.DB
}

func (instance *durability) ParseCredentials(ctx context.Context, raw string) (*entities.Credentials, error) {
	if IsBasicScheme(raw) {
		return ParseBasicCredentials(raw)
	}

	return nil, ErrCredentialsScheme
}

func (instance *durability) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	if err := instance.sequel.Connect(ctx); err != nil {
		return err
	}

	instance.orm = instance.sequel.Client().(*gorm.DB)
	if err := instance.orm.WithContext(ctx).AutoMigrate(&entities.Account{}); err != nil {
		return err
	}

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *durability) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return instance.sequel.Readiness()
}

func (instance *durability) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return instance.sequel.Liveness()
}

func (instance *durability) Disconnect(ctx context.Context) error {
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	instance.status = patterns.StatusDisconnected
	return instance.sequel.Disconnect(ctx)
}

func (instance *durability) Login(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error) {
	if err := entities.ValidateCredentialsOnLogin(credentials); err != nil {
		return nil, err
	}

	var acc entities.Account
	err := instance.orm.WithContext(ctx).
		Where("username = ?", credentials.Username).
		First(&acc).
		Error
	if err != nil {
		return nil, ErrLogin
	}

	if 0 < acc.DeactivatedAt && acc.DeactivatedAt < time.Now().UnixMilli() {
		return nil, ErrAccountDeactivated
	}

	if err := password.CompareString(credentials.Password, acc.PasswordHash); err != nil {
		return nil, ErrLogin
	}

	return acc.Censor(), nil
}

func (instance *durability) Logout(ctx context.Context, credentials *entities.Credentials) error {
	return nil
}

func (instance *durability) Verify(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error) {
	return instance.Login(ctx, credentials)
}

func (instance *durability) Register(ctx context.Context, acc *entities.Account) error {
	if instance.orm.WithContext(ctx).Create(acc).Error != nil {
		return ErrRegister
	}
	return nil
}

func (instance *durability) Deactivate(ctx context.Context, username string, at int64) error {
	return instance.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		acc := &entities.Account{Username: username}

		if tx.First(&acc).Error != nil {
			return ErrAccountNotFound
		}

		// the expired time should be greater than the current one
		if at < acc.DeactivatedAt {
			return ErrDeactivate
		}

		if tx.Model(acc).Update("deactivated_at", at).Error != nil {
			return ErrDeactivate
		}

		return nil
	})
}
