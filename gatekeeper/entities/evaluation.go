package entities

import "github.com/kanthorlabs/common/validator"

type Evaluation struct {
	Tenant   string `json:"tenant" yaml:"tenant"`
	Username string `json:"username" yaml:"username"`
	Role     string `json:"role" yaml:"role"`
}

func (evaluation *Evaluation) Validate() error {
	return validator.Validate(
		validator.StringRequired("OPM.PREVILEGE.TENANT", evaluation.Tenant),
		validator.StringRequired("OPM.PREVILEGE.USERNAME", evaluation.Username),
		validator.StringRequired("OPM.PREVILEGE.ROLE", evaluation.Role),
	)
}
