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

func TestGateway(t *testing.T) {
	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			_, err := New(&config.Config{}, testify.Logger())
			require.ErrorContains(sst, err, "GATEWAY.CONFIG")
		})
	})

	t.Run(".Start/.Stop", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, gw.Start(context.Background()))
		require.ErrorIs(st, gw.Start(context.Background()), ErrAlreadyStarted)
		require.NoError(st, gw.Stop(context.Background()))
		require.ErrorIs(st, gw.Stop(context.Background()), ErrNotStarted)
	})

	t.Run(".UseHttpx/.Run", func(st *testing.T) {
		gw, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, gw.Run(context.Background()), ErrHandlerNotSet)
		require.NoError(st, gw.UseHttpx(chi.NewRouter()))
		require.ErrorIs(st, gw.UseHttpx(chi.NewRouter()), ErrHandlerAlreadySet)
		require.NoError(st, gw.Run(context.Background()))
	})
}
