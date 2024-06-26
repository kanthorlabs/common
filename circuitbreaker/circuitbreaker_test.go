package circuitbreaker

import (
	"testing"

	"github.com/kanthorlabs/common/circuitbreaker/config"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestCircuitBreaker_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := New(testconf, testify.Logger())
		require.NoError(st, err)
	})
}

var testconf = &config.Config{
	Size: 5,
	Close: config.Close{
		CleanupInterval: 5000,
	},
	Half: config.Half{
		PassthroughRequests: 3,
	},
	Open: config.Open{
		Duration: 5000,
		Conditions: config.OpenConditions{
			ErrorConsecutive: 3,
			ErrorRatio:       0.2,
		},
	},
}

func passerror(err error) error {
	return err
}
