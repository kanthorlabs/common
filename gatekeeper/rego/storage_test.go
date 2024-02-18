package rego

import (
	"testing"

	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	t.Run("KO - empty definitions", func(st *testing.T) {
		_, err := Memory(make(map[string][]entities.Permission))
		require.ErrorContains(t, err, "GATEKEEPER.REGO.RBAC.DEFINITION_EMPTY")
	})
}
