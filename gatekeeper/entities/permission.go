package entities

import "github.com/kanthorlabs/common/validator"

var (
	AnyAction = "*"
	AnyObject = "*"
)

type Permission struct {
	Action string `json:"action"`
	Object string `json:"object"`
}

func (permission *Permission) Validate() error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.PERMISSION.ACTION", permission.Action),
		validator.StringRequired("GATEKEEPER.PERMISSION.OBJECT", permission.Object),
	)
}
