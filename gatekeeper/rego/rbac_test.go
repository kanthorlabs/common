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

func TestRBAC(t *testing.T) {
	var rdata regodata
	err := json.Unmarshal(data, &rdata)
	require.Nil(t, err)

	t.Run("New", func(st *testing.T) {
		st.Run("OK", func(sst *testing.T) {
			_, err := RBAC(context.Background(), rdata.Data.Definitions)
			require.Nil(sst, err)
		})

		st.Run("KO - empty definitions", func(sst *testing.T) {
			_, err := RBAC(context.Background(), make(map[string][]entities.Permission))
			require.ErrorContains(sst, err, "GATEKEEPER.REGO.RBAC.DEFINITION_EMPTY")
		})

		st.Run("KO - definition error", func(sst *testing.T) {
			definitions := map[string][]entities.Permission{
				"administrator": {{Action: "*"}},
			}

			_, err = RBAC(context.Background(), definitions)
			require.ErrorContains(sst, err, "GATEKEEPER.PERMISSION.")
		})
	})

	t.Run("Evaluate", func(st *testing.T) {
		evaluate, err := RBAC(context.Background(), rdata.Data.Definitions)
		require.Nil(st, err)

		st.Run("OK", func(sst *testing.T) {
			permission := &entities.Permission{
				Action: "DELETE",
				Object: "/api/application/:id",
			}
			err := evaluate(permission, rdata.Input["administrator"].Privileges)
			require.Nil(sst, err)
		})

		st.Run("KO", func(sst *testing.T) {
			permission := &entities.Permission{
				Action: "DELETE",
				Object: "/api/application/:id",
			}
			err := evaluate(permission, rdata.Input["readonly"].Privileges)
			require.ErrorContains(sst, err, "GATEKEEPER.REGO.RBAC.NOT_ALLOW")
		})
	})
}
