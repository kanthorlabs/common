package sqlx

import (
	"testing"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/stretchr/testify/require"
)

func TestNewGorm(t *testing.T) {
	t.Run("unable to connect", func(st *testing.T) {
		conf := &config.Config{
			Uri: "postgres://postgres:postgres@localhost:2345/postgres",
			Connection: config.Connection{
				MaxLifetime:  300000,
				MaxIdletime:  60000,
				MaxIdleCount: 1,
				MaxOpenCount: 1,
			},
		}
		logger, err := logging.NewNoop()
		require.Nil(t, err)
		_, err = NewGorm(conf, logger)
		require.ErrorContains(t, err, "dial error")
	})
}
