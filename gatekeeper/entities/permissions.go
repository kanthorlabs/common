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
		validator.StringRequired("OPM.PERMISSION.TENANT", permission.Action),
		validator.StringRequired("OPM.PERMISSION.USERNAME", permission.Object),
	)
}
