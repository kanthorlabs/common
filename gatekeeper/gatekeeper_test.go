package gatekeeper

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gatekeeper/config"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGatekeeper_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := New(testconf, testify.Logger())
		require.NoError(st, err)
	})
}

var testconf = &config.Config{
	Engine: config.EngineRBAC,
	Definitions: config.Definitions{
		// Note: decode the base64 value and compare it to the property Definitions at gatekeeper/rego/data.json if you got any error
		Uri: "base64://ewogICJvd25lciI6IFsKICAgIHsKICAgICAgInNjb3BlIjogIioiLAogICAgICAiYWN0aW9uIjogIioiLAogICAgICAib2JqZWN0IjogIioiCiAgICB9CiAgXSwKICAicmVhZG9ubHkiOiBbCiAgICB7CiAgICAgICJzY29wZSI6ICIqIiwKICAgICAgImFjdGlvbiI6ICJHRVQiLAogICAgICAib2JqZWN0IjogIioiCiAgICB9LAogICAgewogICAgICAic2NvcGUiOiAiKiIsCiAgICAgICJhY3Rpb24iOiAiSEVBRCIsCiAgICAgICJvYmplY3QiOiAiKiIKICAgIH0KICBdLAogICJvd24iOiBbCiAgICB7CiAgICAgICJzY29wZSI6ICIqIiwKICAgICAgImFjdGlvbiI6ICIqIiwKICAgICAgIm9iamVjdCI6ICIvYXBpL2FjY291bnQvbWUiCiAgICB9LAogICAgewogICAgICAic2NvcGUiOiAiKiIsCiAgICAgICJhY3Rpb24iOiAiKiIsCiAgICAgICJvYmplY3QiOiAiL2FwaS9hY2NvdW50L3Bhc3N3b3JkIgogICAgfQogIF0sCiAgImFwcGxpY2F0aW9uX3JlYWQiOiBbCiAgICB7CiAgICAgICJzY29wZSI6ICIqIiwKICAgICAgImFjdGlvbiI6ICJHRVQiLAogICAgICAib2JqZWN0IjogIi9hcGkvYXBwbGljYXRpb24iCiAgICB9LAogICAgewogICAgICAic2NvcGUiOiAiKiIsCiAgICAgICJhY3Rpb24iOiAiR0VUIiwKICAgICAgIm9iamVjdCI6ICIvYXBpL2FwcGxpY2F0aW9uLzppZCIKICAgIH0KICBdLAogICJhcHBsaWNhdGlvbl9kZWxldGUiOiBbCiAgICB7CiAgICAgICJzY29wZSI6ICIqIiwKICAgICAgImFjdGlvbiI6ICJERUxFVEUiLAogICAgICAib2JqZWN0IjogIi9hcGkvYXBwbGljYXRpb24vOmlkIgogICAgfQogIF0sCiAgInNjb3BlIjogWwogICAgewogICAgICAic2NvcGUiOiAiaW50ZXJuYWwiLAogICAgICAiYWN0aW9uIjogIlBVVCIsCiAgICAgICJvYmplY3QiOiAiL2FwaS9hcHBsaWNhdGlvbi86aWQiCiAgICB9CiAgXQp9Cg==",
	},
	Privilege: config.Privilege{
		Sqlx: sqlx.Config{
			Uri: testdata.SqliteUri,
			Connection: sqlx.Connection{
				MaxLifetime:  sqlx.DefaultConnMaxLifetime,
				MaxIdletime:  sqlx.DefaultConnMaxIdletime,
				MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
				MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
			},
		},
	},
}

func definitions(t *testing.T) map[string][]entities.Permission {
	definitions, err := config.ParseDefinitionsToPermissions(testconf.Definitions.Uri)
	require.NoError(t, err)

	return definitions
}

func setup(t *testing.T) ([]entities.Privilege, int) {
	privileges := make([]entities.Privilege, 0)

	count := testdata.Fake.IntBetween(5, 10)
	roles := lo.Keys(definitions(t))

	for i := 0; i < count; i++ {
		tenant := uuid.NewString()
		for j := 0; j < count; j++ {
			username := fmt.Sprintf("%s/%s", testdata.Fake.Internet().Email(), uuid.NewString())
			for k := 0; k < len(roles); k++ {
				privilege := entities.Privilege{
					Username:  username,
					Tenant:    tenant,
					Role:      roles[k],
					Metadata:  &safe.Metadata{},
					CreatedAt: time.Now().UnixMilli(),
					UpdatedAt: time.Now().UnixMilli(),
				}
				privileges = append(privileges, privilege)
			}
		}
	}

	return privileges, count
}
