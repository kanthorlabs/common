package rego

import (
	"errors"

	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
)

// Memory is a factory function that returns a storage.Store instance based on the provided definitions.
func Memory(definitions map[string][]entities.Permission) (storage.Store, error) {
	if len(definitions) == 0 {
		return nil, errors.New("GATEKEEPER.REGO.RBAC.DEFINITION_EMPTY.ERROR")
	}

	for role := range definitions {
		for i := range definitions[role] {
			if err := definitions[role][i].Validate(); err != nil {
				return nil, err
			}
		}
	}

	data := map[string]any{
		"permissions": definitions,
	}

	return inmem.NewFromObject(data), nil
}
