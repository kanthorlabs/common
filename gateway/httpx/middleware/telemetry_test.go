package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/opentelemetry"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
)

func TestTelemetry(t *testing.T) {
	opentelemetry.Setup(context.Background())
	t.Cleanup(func() {
		opentelemetry.Teardown(context.Background())
	})

	t.Run("OK", func(st *testing.T) {
		s := chi.NewRouter()
		s.Use(Telemetry())
		s.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			propgator := propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			)
			carrier := propagation.MapCarrier{}
			propgator.Inject(ctx, carrier)

			writer.Ok(w, carrier)
		})

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.NoError(st, err)

		require.NotEmpty(st, body["traceparent"])
	})
}
