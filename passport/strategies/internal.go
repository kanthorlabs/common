package strategies

import (
	"context"
	"errors"
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

func (instance *internal) Register(ctx context.Context, acc entities.Account) error {
	if err := acc.Validate(); err != nil {
		return err
	}
	if err := validator.StringRequired("PASSPORT.ACCOUNT.PASSWORD", acc.Password)(); err != nil {
		return err
	}

	hash, err := password.Hash(acc.Password)
	if err != nil {
		return err
	}

	acc.PasswordHash = hash
	if instance.orm.WithContext(ctx).Create(acc).Error != nil {
		return ErrRegister
	}
	return nil
}

func (instance *internal) Login(ctx context.Context, creds entities.Credentials) (*entities.Tokens, error) {
	return nil, errors.New("PASSPORT.ASK.LOGIN.UNIMPLEMENT.ERROR")
}

func (instance *internal) Logout(ctx context.Context, tokens entities.Tokens) error {
	return errors.New("PASSPORT.ASK.LOGOUT.UNIMPLEMENT.ERROR")
}

func (instance *internal) Verify(ctx context.Context, tokens entities.Tokens) (*entities.Account, error) {
	creds, err := utils.ParseBasicCredentials(tokens.Access)
	if err != nil {
		return nil, err
	}

	if err := creds.Validate(); err != nil {
		return nil, err
	}

	var acc entities.Account
	err = instance.orm.WithContext(ctx).
		Where("username = ?", creds.Username).
		First(&acc).
		Error
	if err != nil {
		return nil, ErrLogin
	}

	if 0 < acc.DeactivatedAt && acc.DeactivatedAt < time.Now().UnixMilli() {
		return nil, ErrAccountDeactivated
	}

	if err := password.Compare(creds.Password, acc.PasswordHash); err != nil {
		return nil, ErrLogin
	}

	return acc.Censor(), nil
}

func (instance *internal) Management() Management {
	if instance.status != patterns.StatusConnected {
		panic(ErrNotConnected)
	}
	return &internalmanagement{orm: instance.orm}
}

type internalmanagement struct {
	orm *gorm.DB
}

func (instance *internalmanagement) Deactivate(ctx context.Context, username string, at int64) error {
	if err := validator.StringRequired("PASSPORT.ACCOUNT.USERNAME", username)(); err != nil {
		return err
	}

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

func (instance *internalmanagement) List(ctx context.Context, usernames []string) ([]*entities.Account, error) {
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
		Order("username DESC").
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

func (instance *internalmanagement) Update(ctx context.Context, acc entities.Account) error {
	err := validator.Validate(
		validator.StringRequired("PASSPORT.ACCOUNT.USERNAME", acc.Username),
		validator.StringRequired("PASSPORT.ACCOUNT.NAME", acc.Name),
		validator.NumberGreaterThan("PASSPORT.ACCOUNT.UPDATED_AT", acc.UpdatedAt, 0),
	)
	if err != nil {
		return err
	}

	return instance.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		acc := &entities.Account{Username: acc.Username}
		if acc.Metadata == nil {
			acc.Metadata = &safe.Metadata{}
		}

		if tx.First(&acc).Error != nil {
			return ErrAccountNotFound
		}

		if acc.DeactivatedAt > 0 && acc.DeactivatedAt < time.Now().UnixMilli() {
			return ErrDeactivate
		}

		updates := map[string]any{"name": acc.Name}
		if acc.Metadata != nil {
			acc.Metadata.Merge(acc.Metadata)
			updates["metadata"] = acc.Metadata
		}

		if tx.Model(acc).Updates(updates).Error != nil {
			return ErrUpdate
		}

		return nil
	})
}
