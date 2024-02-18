package gatekeeper

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gatekeeper/config"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/testdata"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

var testconf = &config.Config{
	Engine: config.EngineRBAC,
	Definitions: config.Definitions{
		// Note: decode the base64 value and compare it to the property Definitions at gatekeeper/rego/data.json if you got any error
		Uri: "base64://ewogICAgICAiYWRtaW5pc3RyYXRvciI6IFsKICAgICAgICB7CiAgICAgICAgICAiYWN0aW9uIjogIioiLAogICAgICAgICAgIm9iamVjdCI6ICIqIgogICAgICAgIH0KICAgICAgXSwKICAgICAgInJlYWRvbmx5IjogWwogICAgICAgIHsKICAgICAgICAgICJhY3Rpb24iOiAiR0VUIiwKICAgICAgICAgICJvYmplY3QiOiAiKiIKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICJhY3Rpb24iOiAiSEVBRCIsCiAgICAgICAgICAib2JqZWN0IjogIioiCiAgICAgICAgfQogICAgICBdLAogICAgICAib3duIjogWwogICAgICAgIHsKICAgICAgICAgICJhY3Rpb24iOiAiKiIsCiAgICAgICAgICAib2JqZWN0IjogIi9hcGkvYWNjb3VudC9tZSIKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICJhY3Rpb24iOiAiKiIsCiAgICAgICAgICAib2JqZWN0IjogIi9hcGkvYWNjb3VudC9wYXNzd29yZCIKICAgICAgICB9CiAgICAgIF0sCiAgICAgICJhcHBsaWNhdGlvbl9yZWFkIjogWwogICAgICAgIHsKICAgICAgICAgICJhY3Rpb24iOiAiR0VUIiwKICAgICAgICAgICJvYmplY3QiOiAiL2FwaS9hcHBsaWNhdGlvbiIKICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICJhY3Rpb24iOiAiR0VUIiwKICAgICAgICAgICJvYmplY3QiOiAiL2FwaS9hcHBsaWNhdGlvbi86aWQiCiAgICAgICAgfQogICAgICBdLAogICAgICAiYXBwbGljYXRpb25fZGVsZXRlIjogWwogICAgICAgIHsKICAgICAgICAgICJhY3Rpb24iOiAiREVMRVRFIiwKICAgICAgICAgICJvYmplY3QiOiAiL2FwaS9hcHBsaWNhdGlvbi86aWQiCiAgICAgICAgfQogICAgICBdCiAgICB9",
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
	require.Nil(t, err)

	return definitions
}

func setup(t *testing.T) ([]entities.Privilege, int) {
	privileges := make([]entities.Privilege, 0)

	count := testdata.Fake.IntBetween(5, 10)
	roles := lo.Keys(definitions(t))

	for i := 0; i < count; i++ {
		tenant := uuid.NewString()
		for j := 0; j < count; j++ {
			username := testdata.Fake.Internet().Email()
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
