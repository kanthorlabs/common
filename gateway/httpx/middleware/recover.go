package middleware

import "github.com/go-chi/chi/v5/middleware"

func Recover() Middleware {
	return middleware.Recoverer
}
