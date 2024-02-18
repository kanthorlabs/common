package entities

import "github.com/kanthorlabs/common/validator"

type Evaluation struct {
	Tenant   string `json:"tenant" yaml:"tenant"`
	Username string `json:"username" yaml:"username"`
	Role     string `json:"role" yaml:"role"`
}

func EvaluationValidateOnGrant(evaluation *Evaluation) error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.EVALUATION.TENANT", evaluation.Tenant),
		validator.StringRequired("GATEKEEPER.EVALUATION.USERNAME", evaluation.Username),
		validator.StringRequired("GATEKEEPER.EVALUATION.ROLE", evaluation.Role),
	)
}

func EvaluationValidateOnRevoke(evaluation *Evaluation) error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.EVALUATION.TENANT", evaluation.Tenant),
		validator.StringRequired("GATEKEEPER.EVALUATION.USERNAME", evaluation.Username),
	)
}

func EvaluationValidateOnEnforce(evaluation *Evaluation) error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.EVALUATION.TENANT", evaluation.Tenant),
		validator.StringRequired("GATEKEEPER.EVALUATION.USERNAME", evaluation.Username),
	)
}
