package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/stretchr/testify/require"
)

func TestRequestId(t *testing.T) {
	s := chi.NewRouter()
	s.Use(RequestId())
	path := "/"
	key := "request_id"
	s.Get(path, func(w http.ResponseWriter, r *http.Request) {
		writer.Ok(w, writer.M{key: r.Context().Value(CtxRequestId)})
	})

	t.Run("OK - use header", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.Nil(st, err)

		id := uuid.NewString()
		req.Header.Set(HeaderRequestId, id)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.Equal(st, id, body[key])
	})

	t.Run("OK - use auto-generate", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.True(st, strings.HasPrefix(body[key].(string), "gw"))
	})
}
