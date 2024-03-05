package config

import (
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/validator"
)

type Ask struct {
	Accounts []entities.Account `json:"accounts" yaml:"accounts" mapstructure:"accounts"`
}

func (conf *Ask) Validate() error {
	return validator.Validate(
		validator.SliceRequired("PASSPORT.STRATEGY.ASK.CONFIG.ACCOUNTS", conf.Accounts),
		validator.Slice(conf.Accounts, func(_ int, acc *entities.Account) error {
			return acc.Validate()
		}),
	)
}
