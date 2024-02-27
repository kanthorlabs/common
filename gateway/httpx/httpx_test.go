package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kanthorlabs/common/gateway/config"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
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

func TestHttpx(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		s, err := New(testconf, testify.Logger())
		require.Nil(st, err)

		path := "/ping"

		s.MethodFunc(http.MethodGet, path, func(w http.ResponseWriter, r *http.Request) {
			writer.Ok(w, writer.M{"pong": true})
		})

		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
	})

	t.Run("OK - panic recover", func(st *testing.T) {
		s, err := New(testconf, testify.Logger())
		require.Nil(st, err)

		path := "/exception"

		s.MethodFunc(http.MethodGet, path, func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("f*ckup"))
		})

		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusInternalServerError, res.Code)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := New(&config.Config{}, testify.Logger())
		require.ErrorContains(st, err, "GATEWAY.CONFIG")
	})
}
