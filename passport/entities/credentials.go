package entities

import "github.com/kanthorlabs/common/validator"

type Credentials struct {
	Region string `json:"region" yaml:"region"`

	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (entity *Credentials) Validate() error {
	return validator.Validate(
		validator.StringRequired("PASSPORT.CREDENTIALS.USERNAME", entity.Username),
		validator.StringRequired("PASSPORT.CREDENTIALS.PASSWORD", entity.Password),
	)
}
