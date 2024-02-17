package rego

import (
	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
)

func memory(definitions map[string][]entities.Permission) (storage.Store, error) {
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
