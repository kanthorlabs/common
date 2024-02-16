package entities

import (
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/timer"
	"github.com/kanthorlabs/common/validator"
)

type Account struct {
	Username     string         `json:"username" yaml:"username" gorm:"primaryKey"`
	PasswordHash string         `json:"password_hash,omitempty" yaml:"password_hash,omitempty"`
	Name         string         `json:"name" yaml:"name"`
	Metadata     *safe.Metadata `json:"metadata" yaml:"metadata"`

	CreatedAt     int64 `json:"created_at" yaml:"created_at"`
	UpdatedAt     int64 `json:"updated_at" yaml:"updated_at"`
	DeactivatedAt int64 `json:"deactivated_at" yaml:"deactivated_at"`
}

func (acc *Account) TableName() string {
	return project.Name("passport_account")
}

func (acc *Account) Validate() error {
	return validator.Validate(
		validator.StringRequired("PASSPORT.ACCOUNT.USERNAME", acc.Username),
		validator.StringRequired("PASSPORT.ACCOUNT.NAME", acc.Name),
		validator.NumberGreaterThan("PASSPORT.ACCOUNT.CREATED_AT", acc.CreatedAt, 0),
		validator.NumberGreaterThan("PASSPORT.ACCOUNT.UPDATED_AT", acc.UpdatedAt, 0),
		validator.NumberGreaterThanOrEqual("PASSPORT.ACCOUNT.DEACTIVATED_AT", acc.DeactivatedAt, 0),
	)
}

func (acc *Account) Censor() *Account {
	censored := &Account{
		Username:  acc.Username,
		Name:      acc.Name,
		Metadata:  &safe.Metadata{},
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
	censored.Metadata.Merge(acc.Metadata)

	return censored
}

func (acc *Account) Active(clock timer.Timer) bool {
	if acc.DeactivatedAt == 0 {
		return true
	}
	return acc.DeactivatedAt > clock.Now().UnixMilli()
}
