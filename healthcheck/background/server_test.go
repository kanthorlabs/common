package background

import (
	"context"
	"testing"

	"github.com/kanthorlabs/common/healthcheck/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestServer_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		name := testdata.Fake.UUID().V4()
		conf := config.Default(name, 1000)
		_, err := NewServer(conf)
		require.NoError(t, err)
	})
	t.Run("KO", func(st *testing.T) {
		_, err := NewServer(&config.Config{})
		require.NotNil(t, err)
	})
}

func TestServer_Readiness(t *testing.T) {
	server, _ := NewServer(&config.Config{
		Dest:      "/",
		Readiness: config.Check{Interval: 1000},
		Liveness:  config.Check{Interval: 1000},
	})
	_ = server.Connect(context.Background())
	defer func() {
		_ = server.Disconnect(context.Background())
	}()

	t.Run("KO - check function error", func(st *testing.T) {
		err := server.Readiness(func() error { return testdata.ErrGeneric })
		require.ErrorIs(st, err, testdata.ErrGeneric)
	})

	t.Run("KO  - write fail", func(st *testing.T) {
		err := server.Readiness(func() error { return nil })
		require.Error(st, err)
	})
}

func TestServer_Liveness(t *testing.T) {
	server, _ := NewServer(&config.Config{
		Dest:      "/",
		Readiness: config.Check{Interval: 1000},
		Liveness:  config.Check{Interval: 1000},
	})
	_ = server.Connect(context.Background())
	defer func() {
		_ = server.Disconnect(context.Background())
	}()

	t.Run("KO - check function error", func(st *testing.T) {
		errc := make(chan error, 1)
		go func() {
			errc <- server.Liveness(func() error { return testdata.ErrGeneric })

		}()
		err := <-errc
		require.ErrorIs(st, err, testdata.ErrGeneric)
	})

	t.Run("KO  - write fail", func(st *testing.T) {
		errc := make(chan error, 1)
		go func() {
			errc <- server.Liveness(func() error { return nil })

		}()
		err := <-errc
		require.Error(st, err)
	})
}
