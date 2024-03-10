package entities

import "github.com/kanthorlabs/common/validator"

var (
	AnyScope  = "*"
	AnyAction = "*"
	AnyObject = "*"
)

type Permission struct {
	Scope  string `json:"scope"`
	Action string `json:"action"`
	Object string `json:"object"`
}

func (permission *Permission) Validate() error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.PERMISSION.SCOPE", permission.Scope),
		validator.StringRequired("GATEKEEPER.PERMISSION.ACTION", permission.Action),
		validator.StringRequired("GATEKEEPER.PERMISSION.OBJECT", permission.Object),
	)
}
