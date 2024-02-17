package entities

import "github.com/kanthorlabs/common/validator"

type Evaluation struct {
	Tenant   string `json:"tenant" yaml:"tenant"`
	Username string `json:"username" yaml:"username"`
	Role     string `json:"role" yaml:"role"`
}

func (evaluation *Evaluation) Validate() error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.PREVILEGE.TENANT", evaluation.Tenant),
		validator.StringRequired("GATEKEEPER.PREVILEGE.USERNAME", evaluation.Username),
		validator.StringRequired("GATEKEEPER.PREVILEGE.ROLE", evaluation.Role),
	)
}
