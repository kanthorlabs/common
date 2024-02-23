package background

import (
	"context"
	"testing"

	"github.com/kanthorlabs/common/healthcheck/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Run("NewServer - KO", func(st *testing.T) {
		_, err := NewServer(&config.Config{})
		require.NotNil(st, err)
	})

	t.Run(".Readiness", func(st *testing.T) {
		server, _ := NewServer(&config.Config{
			Dest:      "/",
			Readiness: config.Check{Timeout: 100, MaxTry: 3},
			Liveness:  config.Check{Timeout: 100, MaxTry: 3},
		})
		_ = server.Connect(context.Background())
		defer func() {
			_ = server.Disconnect(context.Background())
		}()

		st.Run("KO - check function error", func(sst *testing.T) {
			err := server.Readiness(func() error { return testdata.ErrorGeneric })
			require.ErrorIs(sst, err, testdata.ErrorGeneric)
		})

		st.Run("KO  - write fail", func(sst *testing.T) {
			err := server.Readiness(func() error { return nil })
			require.Error(sst, err)
		})
	})

	t.Run(".Liveness", func(st *testing.T) {
		server, _ := NewServer(&config.Config{
			Dest:      "/",
			Readiness: config.Check{Timeout: 100, MaxTry: 3},
			Liveness:  config.Check{Timeout: 100, MaxTry: 3},
		})
		_ = server.Connect(context.Background())
		defer func() {
			_ = server.Disconnect(context.Background())
		}()

		st.Run("KO - check function error", func(sst *testing.T) {
			errc := make(chan error, 1)
			go func() {
				errc <- server.Liveness(func() error { return testdata.ErrorGeneric })

			}()
			err := <-errc
			require.ErrorIs(sst, err, testdata.ErrorGeneric)
		})

		st.Run("KO  - write fail", func(sst *testing.T) {
			errc := make(chan error, 1)
			go func() {
				errc <- server.Liveness(func() error { return nil })

			}()
			err := <-errc
			require.Error(sst, err)
		})
	})
}
