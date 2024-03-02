package background

import (
	"context"
	"testing"

	"github.com/kanthorlabs/common/healthcheck/config"
	"github.com/stretchr/testify/require"
)

func TestBackground(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		// server config
		sconf := config.Default("default", 100)

		server, err := NewServer(sconf)
		require.NoError(t, err)

		go func() {
			ctx := context.Background()
			require.NoError(st, server.Connect(ctx))

			require.NoError(t, server.Readiness(func() error {
				return nil
			}))

			// this function will block the goroutine
			require.NoError(t, server.Liveness(func() error {
				return nil
			}))
		}()

		// client config with default timeout is 10s
		cconf := config.Default("default", 200)
		client, err := NewClient(cconf)
		require.NoError(t, err)

		require.NoError(t, client.Readiness())
		require.NoError(t, client.Liveness())

		require.NoError(st, server.Disconnect(context.Background()))
	})
}
