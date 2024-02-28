package config

import (
	"fmt"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &Config{
			Addr:    fmt.Sprintf(":%d", testdata.Fake.IntBetween(3000, 10000)),
			Timeout: testdata.Fake.Int64Between(1000, 10000),
			Retry: Retry{
				Count:    testdata.Fake.IntBetween(3000, 10000),
				WaitTime: testdata.Fake.Int64Between(1000, 10000),
			},
		}
		require.Nil(st, conf.Validate())
	})

	t.Run("KO ", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(st, conf.Validate(), "SENDER.CONFIG.")
	})

	t.Run("KO - retry error", func(st *testing.T) {
		conf := &Config{
			Addr:    fmt.Sprintf(":%d", testdata.Fake.IntBetween(3000, 10000)),
			Timeout: testdata.Fake.Int64Between(1000, 10000),
		}
		require.ErrorContains(st, conf.Validate(), "SENDER.CONFIG.RETRY")
	})
}
