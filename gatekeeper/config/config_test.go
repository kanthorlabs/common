package config

import (
	"testing"

	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var testconf = &Config{
	Engine: EngineRBAC,
	Definitions: Definitions{
		Uri: "base64://ewogICJhZG1pbmlzdHJhdG9yIjogWwogICAgewogICAgICAiYWN0aW9uIjogIioiLAogICAgICAib2JqZWN0IjogIioiCiAgICB9CiAgXQp9Cg==",
	},
	Privilege: Privilege{
		Sqlx: sqlxconfig.Config{
			Uri: testdata.SqliteUri,
			Connection: sqlxconfig.Connection{
				MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
				MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
				MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
				MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
			},
		},
	},
}

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, testconf.Validate())
	})

	t.Run("KO", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(t, conf.Validate(), "GATEKEEPER.CONFIG.ENGINE")
	})

	t.Run("KO - privilege error", func(st *testing.T) {
		conf := &Config{Engine: EngineRBAC}
		require.ErrorContains(t, conf.Validate(), "SQLX.CONFIG.")
	})

	t.Run("KO - definitions error", func(st *testing.T) {
		conf := &Config{
			Engine:    EngineRBAC,
			Privilege: testconf.Privilege,
		}
		require.ErrorContains(t, conf.Validate(), "GATEKEEPER.CONFIG.DEFINITIONS.")
	})
}
