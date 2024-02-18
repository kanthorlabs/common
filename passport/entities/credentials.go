package entities

import "github.com/kanthorlabs/common/validator"

type Credentials struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func ValidateCredentialsOnLogin(c *Credentials) error {
	if err := validator.PointerNotNil("PASSPORT.CREDENTIALS", c)(); err != nil {
		return err
	}
	return validator.Validate(
		validator.StringRequired("PASSPORT.CREDENTIALS.USERNAME", c.Username),
		validator.StringRequired("PASSPORT.CREDENTIALS.PASSWORD", c.Password),
	)
}
