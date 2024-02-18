package middleware

import (
	"context"
	"net/http"

	"github.com/kanthorlabs/common/idx"
)

var (
	HeaderRequestId        = "X-Request-Id"
	CtxRequestId    ctxkey = "gateway.request.id"
)

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			id := r.Header.Get(HeaderRequestId)
			if id == "" {
				id = idx.New("gw")
			}

			ctx = context.WithValue(ctx, CtxRequestId, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
