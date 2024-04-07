package entities

import "github.com/kanthorlabs/common/validator"

type Tokens struct {
	Access string `json:"access" yaml:"access" mapstructure:"access"`
}

func (entity *Tokens) Validate() error {
	return validator.Validate(
		validator.StringRequired("PASSPORT.TOKENS.ACCESS", entity.Access),
	)
}
