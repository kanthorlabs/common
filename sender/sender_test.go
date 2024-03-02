package sender

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kanthorlabs/common/sender/config"
	"github.com/kanthorlabs/common/sender/entities"
	"github.com/kanthorlabs/common/testify"
	"github.com/kanthorlabs/common/utils"
	"github.com/stretchr/testify/require"
)

var testconf = &config.Config{
	Addr:    ":8080",
	Timeout: 5000,
	Headers: map[string]string{"client": "go-test"},
	Retry: config.Retry{
		Count:    1,
		WaitTime: 500,
	},
}

func TestSender(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := httpserver()
		defer server.Close()

		send, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		id := uuid.NewString()
		req := &entities.Request{
			Method: http.MethodPost,
			Uri:    server.URL + "/200",
			Body:   []byte(fmt.Sprintf(`{"id":"%s"}`, id)),
		}
		res, err := send(context.Background(), req)
		require.NoError(t, err)

		require.Equal(st, http.StatusOK, res.Status)
		require.Contains(st, string(res.Body), id)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := New(&config.Config{}, testify.Logger())
		require.ErrorContains(st, err, "SENDER.CONFIG.")
	})

	t.Run("KO - url parse error", func(st *testing.T) {
		send, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		id := uuid.NewString()
		req := &entities.Request{
			Method: http.MethodPost,
			Uri:    ":://200",
			Body:   []byte(fmt.Sprintf(`{"id":"%s"}`, id)),
		}
		_, err = send(context.Background(), req)
		require.ErrorContains(t, err, "SENDER.URL.PARSE.ERROR")
	})

	t.Run("KO - unsupported scheme error", func(st *testing.T) {
		send, err := New(testconf, testify.Logger())
		require.NoError(st, err)

		id := uuid.NewString()
		req := &entities.Request{
			Method: http.MethodPost,
			Uri:    "/200",
			Body:   []byte(fmt.Sprintf(`{"id":"%s"}`, id)),
		}
		_, err = send(context.Background(), req)
		require.ErrorContains(t, err, "SENDER.SCHEME.NOT_SUPPORT.ERROR")
	})
}

func httpserver() *httptest.Server {
	var handler = func(status int) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)

			body, _ := io.ReadAll(r.Body)
			data := map[string]any{
				"req_headers": r.Header,
				"req_body":    string(body),
			}
			w.Write([]byte(utils.Stringify(data)))
		}
	}

	r := chi.NewRouter()
	r.Get("/200", handler(http.StatusOK))
	r.Post("/200", handler(http.StatusOK))
	r.Put("/200", handler(http.StatusOK))
	r.Patch("/200", handler(http.StatusOK))
	r.Post("/500", handler(http.StatusInternalServerError))
	return httptest.NewServer(r)
}
