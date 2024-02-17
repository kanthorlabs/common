package config

import (
	"testing"

	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run(".Validate", func(st *testing.T) {
		st.Run("KO - engine error", func(sst *testing.T) {
			conf := &Config{}
			require.ErrorContains(sst, conf.Validate(), "GATEKEEPER.CONFIG.ENGINE")
		})

		st.Run("KO - privilege error", func(sst *testing.T) {
			conf := &Config{
				Engine:    EngineRBAC,
				Privilege: Privilege{},
			}
			require.ErrorContains(sst, conf.Validate(), "SQLX.CONFIG.")
		})

		st.Run("OK", func(sst *testing.T) {
			conf := &Config{
				Engine: EngineRBAC,
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
			require.Nil(sst, conf.Validate())
		})
	})
}
