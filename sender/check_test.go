package sender

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/sender/entities"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		mockreq := &entities.Request{
			Method: http.MethodGet,
			Headers: http.Header{
				"Accept":     []string{"*/*; charset=utf-8"},
				"User-Agent": []string{fmt.Sprintf("Kanthor/%s", project.GetVersion())},
			},
			Uri: fmt.Sprintf("http://localhost:%d", testdata.Fake.IntBetween(2000, 10000)),
		}
		mockres := &entities.Response{
			Status:  http.StatusOK,
			Headers: make(http.Header),
			Uri:     mockreq.Uri,
			Body:    []byte("ok"),
		}
		var send = func(ctx context.Context, r *entities.Request) (*entities.Response, error) {
			return mockres, nil
		}

		require.NoError(st, Check(send, mockreq.Uri))
	})

	t.Run("KO - send error", func(st *testing.T) {
		mockreq := &entities.Request{
			Method: http.MethodGet,
			Headers: http.Header{
				"Accept":     []string{"*/*; charset=utf-8"},
				"User-Agent": []string{fmt.Sprintf("Kanthor/%s", project.GetVersion())},
			},
			Uri: fmt.Sprintf("http://localhost:%d", testdata.Fake.IntBetween(2000, 10000)),
		}
		var send = func(ctx context.Context, r *entities.Request) (*entities.Response, error) {
			return nil, testdata.ErrGeneric
		}

		require.ErrorIs(st, Check(send, mockreq.Uri), testdata.ErrGeneric)
	})

	t.Run("KO - status error", func(st *testing.T) {
		mockreq := &entities.Request{
			Method: http.MethodGet,
			Headers: http.Header{
				"Accept":     []string{"*/*; charset=utf-8"},
				"User-Agent": []string{fmt.Sprintf("Kanthor/%s", project.GetVersion())},
			},
			Uri: fmt.Sprintf("http://localhost:%d", testdata.Fake.IntBetween(2000, 10000)),
		}
		mockres := &entities.Response{
			Status:  http.StatusInternalServerError,
			Headers: make(http.Header),
			Uri:     mockreq.Uri,
			Body:    []byte("ko"),
		}
		var send = func(ctx context.Context, r *entities.Request) (*entities.Response, error) {
			return mockres, nil
		}

		require.ErrorContains(st, Check(send, mockreq.Uri), http.StatusText(mockres.Status))
	})
}
