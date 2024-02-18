package rego

import (
	"context"
	_ "embed"
	"errors"

	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/open-policy-agent/opa/rego"
)

//go:embed rbac.rego
var rbac string

func RBAC(ctx context.Context, definitions map[string][]entities.Permission) (Evaluate, error) {
	store, err := Memory(definitions)
	if err != nil {
		return nil, err
	}

	query, err := rego.
		New(
			rego.Query("data.kanthorlabs.gatekeeper.allow"),
			rego.Module("rbac.rego", rbac),
			rego.Store(store),
		).
		PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	return func(permission *entities.Permission, privileges []entities.Privilege) error {
		input := map[string]any{
			"permission": permission,
			"privileges": privileges,
		}
		results, err := query.Eval(ctx, rego.EvalInput(input))
		if err != nil {
			return err
		}

		if !results.Allowed() {
			return errors.New("GATEKEEPER.REGO.RBAC.NOT_ALLOW.ERROR")
		}

		return nil
	}, nil
}
