package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := Connection{
			MaxLifetime:  DefaultConnMaxLifetime,
			MaxIdletime:  DefaultConnMaxIdletime,
			MaxIdleCount: DefaultConnMaxIdleCount,
			MaxOpenCount: DefaultConnMaxOpenCount,
		}
		require.NoError(st, conf.Validate())
	})
	t.Run("KO", func(st *testing.T) {
		conf := Connection{
			MaxLifetime:  0,
			MaxIdletime:  0,
			MaxIdleCount: 0,
			MaxOpenCount: 0,
		}
		require.Error(st, conf.Validate())
	})
}
