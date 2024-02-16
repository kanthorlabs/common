package strategies

import (
	"context"

	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/persistence"
	"github.com/kanthorlabs/common/persistence/sqlx"
	"gorm.io/gorm"
)

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

	orm *gorm.DB
}

func (instance *durability) Connect(ctx context.Context) error {
	if err := instance.sequel.Connect(ctx); err != nil {
		return err
	}

	instance.orm = instance.sequel.Client().(*gorm.DB)
	if err := instance.orm.WithContext(ctx).AutoMigrate(&entities.Account{}); err != nil {
		return err
	}

	return nil
}

func (instance *durability) Readiness() error {
	return instance.sequel.Readiness()
}

func (instance *durability) Liveness() error {
	return instance.sequel.Liveness()
}

func (instance *durability) Disconnect(ctx context.Context) error {
	return instance.sequel.Disconnect(ctx)
}

func (instance *durability) Login(ctx context.Context, credentials *entities.Credentials) (*entities.Account, error) {
	if err := entities.ValidateCredentialsOnLogin(credentials); err != nil {
		return nil, err
	}

	var acc entities.Account
	tx := instance.orm.WithContext(ctx).
		Model(&entities.Account{}).
		Where("username = ?", credentials.Username).
		First(&acc)
	if tx.Error != nil {
		return nil, ErrLogin
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
	tx := instance.orm.WithContext(ctx).Create(acc)
	if tx.Error != nil {
		return ErrRegister
	}
	return nil
}
