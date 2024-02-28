package config

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &Config{
			Size: 5,
			Close: Close{
				CleanupInterval: 5000,
			},
			Half: Half{
				PassthroughRequests: 3,
			},
			Open: Open{
				Duration: 5000,
				Conditions: OpenConditions{
					ErrorConsecutive: 3,
					ErrorRatio:       0.2,
				},
			},
		}
		require.Nil(st, conf.Validate())
	})

	t.Run("KO", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(t, conf.Validate(), "CIRCUIT_BREAKER.CONFIG.SIZE")
	})

	t.Run("KO - close error", func(st *testing.T) {
		conf := &Config{}
		conf.Size = testdata.Fake.IntBetween(1000, 10000)
		require.ErrorContains(t, conf.Validate(), "CIRCUIT_BREAKER.CONFIG.CLOSE")
	})

	t.Run("KO - half error", func(st *testing.T) {
		conf := &Config{}
		conf.Size = testdata.Fake.IntBetween(1000, 10000)
		conf.Close = Close{CleanupInterval: testdata.Fake.Int64Between(1000, 10000)}
		require.ErrorContains(t, conf.Validate(), "CIRCUIT_BREAKER.CONFIG.HALF")
	})

	t.Run("KO - half error", func(st *testing.T) {
		conf := &Config{}
		conf.Size = testdata.Fake.IntBetween(1000, 10000)
		conf.Close = Close{CleanupInterval: testdata.Fake.Int64Between(1000, 10000)}
		require.ErrorContains(t, conf.Validate(), "CIRCUIT_BREAKER.CONFIG.HALF")
	})

	t.Run("KO - open error", func(st *testing.T) {
		conf := &Config{}
		conf.Size = testdata.Fake.IntBetween(1000, 10000)
		conf.Close = Close{CleanupInterval: testdata.Fake.Int64Between(1000, 10000)}
		conf.Half = Half{PassthroughRequests: testdata.Fake.UInt32Between(1000, 10000)}
		require.ErrorContains(t, conf.Validate(), "CIRCUIT_BREAKER.CONFIG.OPEN")
	})

	t.Run("KO - open condition error", func(st *testing.T) {
		conf := &Config{}
		conf.Size = testdata.Fake.IntBetween(1000, 10000)
		conf.Close = Close{CleanupInterval: testdata.Fake.Int64Between(1000, 10000)}
		conf.Half = Half{PassthroughRequests: testdata.Fake.UInt32Between(1000, 10000)}
		conf.Open = Open{Duration: testdata.Fake.Int64Between(1000, 10000)}
		require.ErrorContains(t, conf.Validate(), "CIRCUIT_BREAKER.CONFIG.OPEN.CONDITION")
	})
}
