package middleware

import (
	"context"
	"net/http"
	"slices"

	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/idempotency"
)

var HeaderIdempotencyKey = "Idempotency-Key"
var requires = []string{
	http.MethodPost,
}

func Idempotency(idemp idempotency.Idempotency, bypass bool) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if bypass {
				next.ServeHTTP(w, r)
				return
			}

			key := r.Header.Get(HeaderIdempotencyKey)
			if key == "" {
				writer.ErrBadRequest(w, writer.ErrorString("GATEWAY.IDEMPOTENCY.KEY_EMPTY.ERROR"))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, idempotency.CtxKey, key)

			if slices.Contains(requires, r.Method) {
				if err := idemp.Validate(ctx, key); err != nil {
					writer.ErrConflict(w, writer.Error(err))
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
