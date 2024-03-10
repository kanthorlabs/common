package config

import (
	"fmt"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("Ok", func(st *testing.T) {
		conf := &Config{
			Addr:    ":8080",
			Timeout: 60000,
			Cors: Cors{
				MaxAge: 86400000,
			}}
		require.NoError(st, conf.Validate())
	})

	t.Run("KO", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(st, conf.Validate(), "GATEWAY.CONFIG.")
	})

	t.Run("KO - CORS error", func(st *testing.T) {
		conf := &Config{
			Addr:    fmt.Sprintf(":%d", testdata.Fake.IntBetween(3000, 10000)),
			Timeout: testdata.Fake.Int64Between(1000, 10000),
		}
		require.ErrorContains(st, conf.Validate(), "GATEWAY.CONFIG.CORS")
	})
}
