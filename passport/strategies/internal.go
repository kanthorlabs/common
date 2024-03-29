package strategies

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/passport/utils"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/sqlx"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/validator"
	"gorm.io/gorm"
)

// NewInternal creates a new internal strategy instance what allows to authenticate users based on a username and password.
// The storage is based on a SQL database.
func NewInternal(conf *config.Internal, logger logging.Logger) (Strategy, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	sequel, err := sqlx.New(&conf.Sqlx, logger)
	if err != nil {
		return nil, err
	}

	return &internal{conf: conf, logger: logger, sequel: sequel}, nil
}

type internal struct {
	conf   *config.Internal
	logger logging.Logger
	sequel persistence.Persistence

	mu     sync.Mutex
	status int
	orm    *gorm.DB
}

func (instance *internal) ParseCredentials(ctx context.Context, raw string) (*entities.Credentials, error) {
	if utils.IsBasicScheme(raw) {
		creds, err := utils.ParseBasicCredentials(raw)
		if err != nil {
			instance.logger.Error(err.Error())
			return nil, ErrParseCredentials
		}

		return creds, nil
	}

	return nil, ErrCredentialsScheme
}

func (instance *internal) Connect(ctx context.Context) error {
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

func (instance *internal) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return instance.sequel.Readiness()
}

func (instance *internal) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return instance.sequel.Liveness()
}

func (instance *internal) Disconnect(ctx context.Context) error {
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	instance.status = patterns.StatusDisconnected
	return instance.sequel.Disconnect(ctx)
}

func (instance *internal) Login(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error) {
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

	if err := password.Compare(credentials.Password, acc.PasswordHash); err != nil {
		return nil, ErrLogin
	}

	return acc.Censor(), nil
}

func (instance *internal) Logout(ctx context.Context, credentials *entities.Credentials) error {
	return nil
}

func (instance *internal) Verify(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error) {
	return instance.Login(ctx, credentials)
}

func (instance *internal) Register(ctx context.Context, acc *entities.Account) error {
	if instance.orm.WithContext(ctx).Create(acc).Error != nil {
		return ErrRegister
	}
	return nil
}

func (instance *internal) Deactivate(ctx context.Context, username string, at int64) error {
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

func (instance *internal) List(ctx context.Context, usernames []string) ([]*entities.Account, error) {
	err := validator.Validate(
		validator.SliceRequired("usernames", usernames),
		validator.Slice(usernames, func(i int, item *string) error {
			key := fmt.Sprintf("usernames[%d]", i)
			return validator.StringRequired(key, *item)()
		}),
	)
	if err != nil {
		return nil, err
	}

	var docs []entities.Account
	err = instance.orm.WithContext(ctx).
		Where("username IN ?", usernames).
		Find(&docs).
		Error
	if err != nil {
		return nil, ErrList
	}

	accounts := make([]*entities.Account, len(docs))
	for i := range docs {
		accounts[i] = docs[i].Censor()
	}

	return accounts, nil
}

func (instance *internal) Update(ctx context.Context, account *entities.Account) error {
	return instance.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		acc := &entities.Account{Username: account.Username}
		if acc.Metadata == nil {
			acc.Metadata = &safe.Metadata{}
		}

		if tx.First(&acc).Error != nil {
			return ErrAccountNotFound
		}

		if acc.DeactivatedAt > 0 && acc.DeactivatedAt < time.Now().UnixMilli() {
			return ErrDeactivate
		}

		updates := map[string]any{}
		if account.Name != "" {
			updates["name"] = account.Name
		}
		if account.Metadata != nil {
			account.Metadata.Merge(acc.Metadata)
			updates["metadata"] = account.Metadata
		}

		if tx.Model(acc).Updates(updates).Error != nil {
			return ErrUpdate
		}

		return nil
	})
}
