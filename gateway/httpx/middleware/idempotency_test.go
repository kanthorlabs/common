package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/idempotency"
	mocidemp "github.com/kanthorlabs/common/mocks/idempotency"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestIdempotency(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		s := chi.NewRouter()
		idemp := mocidemp.NewIdempotency(t)
		s.Use(Idempotency(idemp, false))

		path := "/" + testdata.Fake.RandomStringWithLength(10)
		s.Post(path, func(w http.ResponseWriter, r *http.Request) {
			key := r.Context().Value(idempotency.CtxKey).(string)
			writer.Ok(w, writer.M{"key": key})
		})

		req, err := http.NewRequest(http.MethodPost, path, nil)
		require.Nil(st, err)

		key := uuid.NewString()
		req.Header.Set(HeaderIdempotencyKey, key)

		// only care about the engine
		idemp.EXPECT().
			Validate(mock.Anything, key).
			Return(nil).
			Once()

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.Contains(st, body["key"], key)
	})

	t.Run("OK - bypass", func(st *testing.T) {
		s := chi.NewRouter()
		idemp := mocidemp.NewIdempotency(t)
		s.Use(Idempotency(idemp, true))

		path := "/" + testdata.Fake.RandomStringWithLength(10)
		s.Get(path, func(w http.ResponseWriter, r *http.Request) {
			key, ok := r.Context().Value(idempotency.CtxKey).(string)
			writer.Ok(w, writer.M{"key": key, "ok": ok})
		})

		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.Nil(st, err)

		key := uuid.NewString()
		req.Header.Set(HeaderIdempotencyKey, key)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.False(st, body["ok"].(bool))
	})

	t.Run("KO - empty key error", func(st *testing.T) {
		s := chi.NewRouter()
		idemp := mocidemp.NewIdempotency(t)
		s.Use(Idempotency(idemp, false))

		path := "/" + testdata.Fake.RandomStringWithLength(10)
		s.Post(path, func(w http.ResponseWriter, r *http.Request) {
			key := r.Context().Value(idempotency.CtxKey).(string)
			writer.Ok(w, writer.M{"key": key})
		})

		req, err := http.NewRequest(http.MethodPost, path, nil)
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusBadRequest, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.Contains(st, body["error"], "GATEWAY.IDEMPOTENCY.KEY_EMPTY.ERROR")
	})

	t.Run("KO - already seen error", func(st *testing.T) {
		s := chi.NewRouter()
		idemp := mocidemp.NewIdempotency(t)
		s.Use(Idempotency(idemp, false))

		path := "/" + testdata.Fake.RandomStringWithLength(10)
		s.Post(path, func(w http.ResponseWriter, r *http.Request) {
			key := r.Context().Value(idempotency.CtxKey).(string)
			writer.Ok(w, writer.M{"key": key})
		})

		req, err := http.NewRequest(http.MethodPost, path, nil)
		require.Nil(st, err)

		key := uuid.NewString()
		req.Header.Set(HeaderIdempotencyKey, key)

		// only care about the engine
		idemp.EXPECT().
			Validate(mock.Anything, key).
			Return(idempotency.ErrConflict).
			Once()

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusConflict, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.Contains(st, body["error"], idempotency.ErrConflict.Error())
	})

}
