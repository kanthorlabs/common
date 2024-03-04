package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestUseHealthz(t *testing.T) {
	s, err := New(testconf, testify.Logger())
	require.NoError(t, err)

	ok := "/helthz/ok"
	ko := "/healthz/ko"
	s.Get(ok, UseHealthz(func() error { return nil }))
	s.Get(ko, UseHealthz(func() error { return testdata.ErrGeneric }))

	t.Run("OK", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ok, nil)
		require.NoError(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
	})

	t.Run("KO", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ko, nil)
		require.NoError(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusInternalServerError, res.Code)
	})
}

func TestUseVersion(t *testing.T) {
	s, err := New(testconf, testify.Logger())
	require.NoError(t, err)

	s.Get("/", UseHealthz(func() error { return nil }))

	t.Run("OK", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
	})
}
