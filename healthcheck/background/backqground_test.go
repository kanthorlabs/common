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
		require.Nil(t, err)

		go func() {
			ctx := context.Background()
			require.Nil(st, server.Connect(ctx))

			require.Nil(t, server.Readiness(func() error {
				return nil
			}))

			// this function will block the goroutine
			require.Nil(t, server.Liveness(func() error {
				return nil
			}))
		}()

		// client config with default timeout is 10s
		cconf := config.Default("default", 200)
		client, err := NewClient(cconf)
		require.Nil(t, err)

		require.Nil(t, client.Readiness())
		require.Nil(t, client.Liveness())

		require.Nil(st, server.Disconnect(context.Background()))
	})
}
