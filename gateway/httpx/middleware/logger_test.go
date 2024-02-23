package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/mocks/logging"
	"github.com/kanthorlabs/common/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("OK - 2xx", func(st *testing.T) {
		logger := logging.NewLogger(st)
		logger.EXPECT().Debugw(
			mock.AnythingOfType("string"), // message
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("int"),    // status
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // method
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // uri
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // headers
		).Return()

		s := chi.NewRouter()
		s.Use(Logger(logger))
		path := "/"
		s.Get(path, func(w http.ResponseWriter, r *http.Request) {
			writer.Ok(w, writer.M{})
		})

		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
		logger.AssertCalled(
			st, "Debugw",
			"GATEWAY.REQUEST",
			"response_status", http.StatusOK,
			"request_method", http.MethodGet,
			"request_uri", path,
			"request_headers", utils.Stringify(map[string]any{}),
		)
	})

	t.Run("OK - GET 5xx", func(st *testing.T) {
		logger := logging.NewLogger(st)
		logger.EXPECT().Errorw(
			mock.AnythingOfType("string"), // message
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("int"),    // status
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // method
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // uri
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // headers
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // body
		).Return()

		s := chi.NewRouter()
		s.Use(Logger(logger))
		path := "/500"

		s.Get(path, func(w http.ResponseWriter, r *http.Request) {
			writer.ErrUnknown(w, writer.M{})
		})

		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusInternalServerError, res.Code)
		logger.AssertCalled(
			st, "Errorw",
			"GATEWAY.REQUEST.ERROR",
			"response_status", http.StatusInternalServerError,
			"request_method", http.MethodGet,
			"request_uri", path,
			"request_headers", utils.Stringify(map[string]any{}),
			"request_body", "",
		)
	})

	t.Run("OK - POST 5xx", func(st *testing.T) {
		logger := logging.NewLogger(st)
		logger.EXPECT().Errorw(
			mock.AnythingOfType("string"), // message
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("int"),    // status
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // method
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // uri
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // headers
			mock.AnythingOfType("string"), // key
			mock.AnythingOfType("string"), // body
		).Return()

		s := chi.NewRouter()
		s.Use(Logger(logger))
		path := "/500"

		s.Post(path, func(w http.ResponseWriter, r *http.Request) {
			writer.ErrUnknown(w, writer.M{})
		})

		body := utils.Stringify(map[string]string{"id": uuid.NewString()})
		req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(body))
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusInternalServerError, res.Code)
		logger.AssertCalled(
			st, "Errorw",
			"GATEWAY.REQUEST.ERROR",
			"response_status", http.StatusInternalServerError,
			"request_method", http.MethodPost,
			"request_uri", path,
			"request_headers", utils.Stringify(map[string]any{}),
			"request_body", body,
		)
	})
}
