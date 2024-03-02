package config

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gatekeeper/entities"
	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/utils"
	"github.com/stretchr/testify/require"
)

var testconf = &Config{
	Engine: EngineRBAC,
	Definitions: Definitions{
		Uri: "base64://ewogICJhZG1pbmlzdHJhdG9yIjogWwogICAgewogICAgICAiYWN0aW9uIjogIioiLAogICAgICAib2JqZWN0IjogIioiCiAgICB9CiAgXQp9Cg==",
	},
	Privilege: Privilege{
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

func TestConfig(t *testing.T) {
	t.Run(".Validate", func(st *testing.T) {
		st.Run("KO - engine error", func(sst *testing.T) {
			conf := &Config{}
			require.ErrorContains(sst, conf.Validate(), "GATEKEEPER.CONFIG.ENGINE")
		})

		st.Run("KO - privilege error", func(sst *testing.T) {
			conf := &Config{Engine: EngineRBAC}
			require.ErrorContains(sst, conf.Validate(), "SQLX.CONFIG.")
		})

		st.Run("KO - definitions error", func(sst *testing.T) {
			conf := &Config{
				Engine:    EngineRBAC,
				Privilege: testconf.Privilege,
			}
			require.ErrorContains(sst, conf.Validate(), "GATEKEEPER.CONFIG.DEFINITIONS.")
		})

		st.Run("OK", func(sst *testing.T) {
			require.NoError(sst, testconf.Validate())
		})
	})
}

func TestParseDefinitionsToPermissions(t *testing.T) {
	definitions := map[string][]entities.Permission{
		"administrator": {{Action: "*", Object: "*"}},
	}

	permission := utils.Stringify(definitions)

	dpath := os.TempDir() + "/" + uuid.NewString()
	require.NoError(t, os.WriteFile(dpath, []byte(permission), os.ModePerm))
	dbase64 := base64.StdEncoding.EncodeToString([]byte(permission))

	t.Run("file", func(st *testing.T) {
		st.Run("OK", func(sst *testing.T) {
			definitions, err := ParseDefinitionsToPermissions("file://" + dpath)
			require.NoError(t, err)

			require.Equal(sst, definitions, definitions)
		})

		st.Run("KO - file not found", func(sst *testing.T) {
			_, err := ParseDefinitionsToPermissions("file://./not-found/file")
			require.ErrorContains(st, err, "no such file or directory")
		})

		st.Run("KO - unmarshal error", func(sst *testing.T) {
			p := os.TempDir() + "/" + uuid.NewString()
			require.NoError(sst, os.WriteFile(p, []byte(""), os.ModePerm))

			_, err := ParseDefinitionsToPermissions("file://" + p)
			require.ErrorContains(st, err, "unexpected end of JSON input")
		})

	})

	t.Run("base64", func(st *testing.T) {
		st.Run("OK", func(sst *testing.T) {
			definitions, err := ParseDefinitionsToPermissions("base64://" + dbase64)
			require.NoError(t, err)

			require.Equal(sst, definitions, definitions)
		})

		st.Run("KO - decode error", func(sst *testing.T) {
			_, err := ParseDefinitionsToPermissions("base64://---")
			require.Equal(t, err, base64.CorruptInputError(0))
		})

		st.Run("KO - unmarshal error", func(sst *testing.T) {
			_, err := ParseDefinitionsToPermissions("base64://ey19")
			require.ErrorContains(t, err, "invalid character")
		})
	})

	t.Run("KO - unsupported uri", func(st *testing.T) {
		_, err := ParseDefinitionsToPermissions(testdata.Fake.Internet().URL())
		require.ErrorContains(st, err, "GATEKEEPER.CONFIG.DEFINITIONS.URI.UNSUPPORTED")
	})
}
