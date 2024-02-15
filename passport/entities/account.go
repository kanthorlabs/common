package entities

import (
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/validator"
)

var (
	TenantSuper = "*"
)

type Account struct {
	Username     string         `json:"username" yaml:"username"`
	PasswordHash string         `json:"password_hash,omitempty" yaml:"password_hash,omitempty"`
	Tenant       string         `json:"tenant" yaml:"tenant"`
	Name         string         `json:"name" yaml:"name"`
	Metadata     *safe.Metadata `json:"metadata" yaml:"metadata"`

	CreatedAt int64 `json:"created_at" yaml:"created_at"`
	UpdatedAt int64 `json:"updated_at" yaml:"updated_at"`
}

func (acc *Account) Validate() error {
	return validator.Validate(
		validator.StringRequired("PASSPORT.ACCOUNT.USERNAME", acc.Username),
		validator.StringRequired("PASSPORT.ACCOUNT.PASSWORD_HASH", acc.PasswordHash),
		validator.StringRequired("PASSPORT.ACCOUNT.NAME", acc.Name),
		validator.StringRequired("PASSPORT.ACCOUNT.TENANT", acc.Tenant),
		validator.NumberGreaterThan("PASSPORT.ACCOUNT.CREATED_AT", acc.CreatedAt, 0),
		validator.NumberGreaterThan("PASSPORT.ACCOUNT.UPDATED_AT", acc.UpdatedAt, 0),
	)
}

func (acc *Account) Censor() *Account {
	censored := &Account{
		Username:  acc.Username,
		Tenant:    acc.Tenant,
		Name:      acc.Name,
		Metadata:  &safe.Metadata{},
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
	censored.Metadata.Merge(acc.Metadata)

	return censored
}
