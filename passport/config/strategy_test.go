package config

import (
	"testing"

	"github.com/google/uuid"
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestStrategy(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &Strategy{
			Engine: EngineDurability,
			Name:   uuid.NewString(),
			Durability: Durability{
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
		require.NoError(st, conf.Validate())
	})

	t.Run("KO - enginee error", func(st *testing.T) {
		conf := &Strategy{}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.CONFIG.NAME")
	})

	t.Run("KO - Ask error", func(st *testing.T) {
		conf := &Strategy{Engine: EngineAsk, Name: uuid.NewString()}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.ASK.CONFIG")
	})

	t.Run("KO - Durability error", func(st *testing.T) {
		conf := &Strategy{Engine: EngineDurability, Name: uuid.NewString()}
		require.ErrorContains(st, conf.Validate(), "SQLX.CONFIG.")
	})
}
