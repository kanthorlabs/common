package sender

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/sender/entities"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestHttp(t *testing.T) {
	send, err := Http(testconf, testify.Logger())
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		server := httpserver()
		defer server.Close()

		methods := map[string]bool{
			http.MethodGet:   false,
			http.MethodPost:  true,
			http.MethodPut:   true,
			http.MethodPatch: true,
		}

		for method, hasBody := range methods {
			id := uuid.NewString()
			req := &entities.Request{
				Method: method,
				Uri:    server.URL + "/200",
				Body:   []byte(fmt.Sprintf(`{"id":"%s"}`, id)),
			}
			res, err := send(context.Background(), req)
			require.NoError(t, err)

			require.Equal(st, http.StatusOK, res.Status)

			if hasBody {
				require.Contains(st, string(res.Body), id)
			} else {
				require.NotContains(st, string(res.Body), id)
			}
		}
	})

	t.Run("OK - retry", func(st *testing.T) {
		server := httpserver()
		defer server.Close()

		id := uuid.NewString()
		req := &entities.Request{
			Method: http.MethodPost,
			Uri:    server.URL + "/500",
			Body:   []byte(fmt.Sprintf(`{"id":"%s"}`, id)),
		}
		res, err := send(context.Background(), req)
		require.NoError(t, err)

		require.Equal(st, http.StatusInternalServerError, res.Status)
	})

	t.Run("KO - request validation error", func(st *testing.T) {
		id := uuid.NewString()
		req := &entities.Request{
			Method: http.MethodPost,
			Body:   []byte(fmt.Sprintf(`{"id":"%s"}`, id)),
		}
		_, err := send(context.Background(), req)
		require.ErrorContains(t, err, "SENDER.REQUEST")
	})

	t.Run("KO - timeout", func(st *testing.T) {
		server := httpserver()
		defer server.Close()

		id := uuid.NewString()
		req := &entities.Request{
			Method: http.MethodPost,
			Uri:    fmt.Sprintf("%s/delay?ms=%d", server.URL, testconf.Timeout),
			Body:   []byte(fmt.Sprintf(`{"id":"%s"}`, id)),
		}
		res, err := send(context.Background(), req)
		require.NoError(t, err)

		require.Equal(st, -1, res.Status)
		require.Empty(st, res.Headers)
		require.NotEmpty(st, res.Body)
	})
}
