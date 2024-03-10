package rego

import (
	"testing"

	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := Memory(map[string][]entities.Permission{
			"administrator": {{
				Scope:  entities.AnyScope,
				Action: entities.AnyAction,
				Object: entities.AnyObject,
			}},
		})
		require.NoError(t, err)
	})

	t.Run("KO - empty definitions", func(st *testing.T) {
		_, err := Memory(make(map[string][]entities.Permission))
		require.ErrorContains(t, err, "GATEKEEPER.REGO.RBAC.DEFINITION_EMPTY")
	})

	t.Run("KO - definition error", func(st *testing.T) {
		_, err := Memory(map[string][]entities.Permission{
			"administrator": {{Scope: entities.AnyScope}},
		})
		require.ErrorContains(t, err, "GATEKEEPER.PERMISSION.ACTION")
	})
}
