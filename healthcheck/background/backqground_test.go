package background

import (
	"context"
	"testing"

	"github.com/kanthorlabs/common/healthcheck/config"
	"github.com/stretchr/testify/require"
)

func TestBackground_Readiness(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		// server config
		sconf := config.Default("default", 1000)

		server, err := NewServer(sconf)
		require.NoError(t, err)

		done := make(chan bool, 1)
		go func() {
			ctx := context.Background()
			require.NoError(st, server.Connect(ctx))

			require.NoError(t, server.Readiness(func() error {
				return nil
			}))

			done <- true
		}()

		<-done
		cconf := config.Default("default", 2000)
		client, err := NewClient(cconf)
		require.NoError(st, err)

		require.NoError(st, client.Readiness())
		require.NoError(st, server.Disconnect(context.Background()))
	})
}

func TestBackground_Liveness(t *testing.T) {
	t.Run("default", func(st *testing.T) {
		// server config
		sconf := config.Default("default", 1000)

		server, err := NewServer(sconf)
		require.NoError(t, err)

		done := make(chan bool, 1)
		go func() {
			ctx := context.Background()
			require.NoError(st, server.Connect(ctx))
			var ok bool

			require.NoError(t, server.Liveness(func() error {
				if ok {
					done <- ok
				}

				ok = true
				return nil
			}))

		}()

		<-done
		cconf := config.Default("default", 2000)
		client, err := NewClient(cconf)
		require.NoError(st, err)

		require.NoError(st, client.Liveness())
		require.NoError(st, server.Disconnect(context.Background()))
	})
}
