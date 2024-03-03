package gateway

import (
	"context"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gateway/config"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

var testconf = &config.Config{
	Addr:    ":8080",
	Timeout: 60000,
	Cors: config.Cors{
		MaxAge: 86400000,
	},
}

func TestGateway_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := New(testconf, testify.Logger())
		require.NoError(t, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := New(&config.Config{}, testify.Logger())
		require.ErrorContains(t, err, "GATEWAY.CONFIG")
	})
}

func TestGateway_Start(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, gw.Start(context.Background()))
	})

	t.Run("KO - already started", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, gw.Start(context.Background()))
		require.ErrorIs(t, gw.Start(context.Background()), ErrAlreadyStarted)
	})
}

func TestGateway_Stop(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(t, err)

		gw.UseHttpx(chi.NewRouter())
		require.NoError(t, gw.Start(context.Background()))
		require.NoError(t, gw.Stop(context.Background()))
	})

	t.Run("KO - not started", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, gw.Stop(context.Background()), ErrNotStarted)
	})
}

func TestGateway_Run(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(t, err)

		gw.UseHttpx(chi.NewRouter())
		require.NoError(t, gw.Start(context.Background()))
		require.NoError(t, gw.Run(context.Background()))
		require.NoError(t, gw.Stop(context.Background()))
	})

	t.Run("KO - handler not set", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(t, err)

		require.ErrorIs(t, gw.Run(context.Background()), ErrHandlerNotSet)
	})
}

func TestGateway_UseHttpx(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, gw.UseHttpx(chi.NewRouter()))
	})

	t.Run("KO - handler already set", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(t, err)

		require.NoError(t, gw.UseHttpx(chi.NewRouter()))
		require.ErrorIs(t, gw.UseHttpx(chi.NewRouter()), ErrHandlerAlreadySet)
	})
}
