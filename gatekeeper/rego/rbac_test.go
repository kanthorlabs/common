package rego

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/stretchr/testify/require"
)

//go:embed data.json
var data []byte

type regodata struct {
	Data struct {
		Definitions map[string][]entities.Permission
	}
	Input map[string]struct {
		Privileges []entities.Privilege
	}
}

func TestRBAC_New(t *testing.T) {
	var rdata regodata
	err := json.Unmarshal(data, &rdata)
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		_, err := RBAC(context.Background(), rdata.Data.Definitions)
		require.NoError(st, err)
	})

	t.Run("KO - empty definitions", func(st *testing.T) {
		_, err := RBAC(context.Background(), make(map[string][]entities.Permission))
		require.ErrorContains(st, err, "GATEKEEPER.REGO.RBAC.DEFINITION_EMPTY")
	})

	t.Run("KO - definition error", func(st *testing.T) {
		definitions := map[string][]entities.Permission{
			"administrator": {{Action: "*"}},
		}

		_, err = RBAC(context.Background(), definitions)
		require.ErrorContains(st, err, "GATEKEEPER.PERMISSION.")
	})
}

func TestRBAC_Evaluate(t *testing.T) {
	var rdata regodata
	err := json.Unmarshal(data, &rdata)
	require.NoError(t, err)

	evaluate, err := RBAC(context.Background(), rdata.Data.Definitions)
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		permission := &entities.Permission{
			Action: "DELETE",
			Object: "/api/application/:id",
		}
		err := evaluate(permission, rdata.Input["administrator"].Privileges)
		require.NoError(st, err)
	})

	t.Run("KO", func(st *testing.T) {
		permission := &entities.Permission{
			Action: "DELETE",
			Object: "/api/application/:id",
		}
		err := evaluate(permission, rdata.Input["readonly"].Privileges)
		require.ErrorContains(st, err, "GATEKEEPER.REGO.RBAC.NOT_ALLOW")
	})
}
