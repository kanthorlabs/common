package writer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	s := chi.NewRouter()

	key := "request_id"
	id := uuid.NewString()
	reply := M{key: id}

	s.Post("/200", func(w http.ResponseWriter, r *http.Request) {
		Ok(w, reply)
	})
	s.Post("/201", func(w http.ResponseWriter, r *http.Request) {
		Created(w, reply)
	})
	s.Post("/400", func(w http.ResponseWriter, r *http.Request) {
		ErrBadRequest(w, reply)
	})
	s.Post("/401", func(w http.ResponseWriter, r *http.Request) {
		ErrUnauthorized(w, reply)
	})
	s.Post("/404", func(w http.ResponseWriter, r *http.Request) {
		ErrNotFound(w, reply)
	})
	s.Post("/409", func(w http.ResponseWriter, r *http.Request) {
		ErrConflict(w, reply)
	})
	s.Post("/500", func(w http.ResponseWriter, r *http.Request) {
		ErrUnknown(w, reply)
	})

	status := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusConflict,
		http.StatusInternalServerError,
	}

	for i := range status {
		path := fmt.Sprintf("/%d", status[i])
		req, err := http.NewRequest(http.MethodPost, path, nil)
		require.Nil(t, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(t, status[i], res.Code)

		var body M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(t, err)

		require.Contains(t, reply[key], body[key])
	}

	msg := testdata.Fake.Lorem().Sentence(10)

	s.Post("/error", func(w http.ResponseWriter, r *http.Request) {
		ErrUnknown(w, Error(errors.New(msg)))
	})

	s.Post("/error/string", func(w http.ResponseWriter, r *http.Request) {
		ErrUnknown(w, ErrorString(msg))
	})

	errs := []string{"/error", "/error/string"}
	for i := range errs {
		path := errs[i]
		req, err := http.NewRequest(http.MethodPost, path, nil)
		require.Nil(t, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(t, http.StatusInternalServerError, res.Code)

		var body M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(t, err)

		require.Contains(t, msg, body["error"])
	}

}
