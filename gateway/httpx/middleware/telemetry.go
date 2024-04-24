package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

func Telemetry() Middleware {
	tracer := otel.Tracer("GATEWAY")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.URL.String()))
			defer span.End()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
