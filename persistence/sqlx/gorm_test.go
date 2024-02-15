package sqlx

import (
	"testing"

	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestNewGorm(t *testing.T) {
	t.Run("KO because of connection error", func(st *testing.T) {
		conf := &config.Config{
			Uri: "postgres://postgres:postgres@localhost:2345/postgres",
			Connection: config.Connection{
				MaxLifetime:  300000,
				MaxIdletime:  60000,
				MaxIdleCount: 1,
				MaxOpenCount: 1,
			},
		}
		_, err := NewGorm(conf, testify.Logger())
		require.ErrorContains(t, err, "dial error")
	})
}
