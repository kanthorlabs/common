package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestRecover(t *testing.T) {
	s := chi.NewRouter()
	s.Use(Recover())
	path := "/"
	s.Get(path, func(w http.ResponseWriter, r *http.Request) {
		panic(testdata.ErrGeneric)
	})

	req, err := http.NewRequest(http.MethodGet, path, nil)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	s.ServeHTTP(res, req)

	require.Equal(t, http.StatusInternalServerError, res.Code)
}
