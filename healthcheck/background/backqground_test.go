package background

import (
	"context"
	"fmt"
	"testing"
	"time"

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

		ctx := context.Background()
		require.NoError(st, server.Connect(ctx))

		// liveness in server side locks goroutine
		go server.Liveness(func() error { return nil })

		cconf := config.Default("default", 2000)
		client, err := NewClient(cconf)
		require.NoError(st, err)

		tries := 10
		for i := 0; i < tries; i++ {
			if err := client.Liveness(); err == nil {
				require.NoError(st, server.Disconnect(context.Background()))
				return
			}
			time.Sleep(time.Second * time.Duration((i + 1)))
		}

		require.NoError(st, server.Disconnect(context.Background()))
		panic(fmt.Errorf("failed check after %d tries", tries))
	})
}
