package httpx

import (
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kanthorlabs/common/gateway/config"
	"github.com/kanthorlabs/common/gateway/httpx/middleware"
	"github.com/kanthorlabs/common/logging"
)

type Httpx interface {
	chi.Router
}

func New(conf *config.Config, logger logging.Logger) (Httpx, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recover())
	r.Use(chimw.Timeout(time.Millisecond * time.Duration(conf.Timeout)))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   conf.Cors.AllowedOrigins,
		AllowedMethods:   conf.Cors.AllowedMethods,
		AllowedHeaders:   conf.Cors.AllowedHeaders,
		ExposedHeaders:   conf.Cors.ExposedHeaders,
		AllowCredentials: conf.Cors.AllowCredentials,
		MaxAge:           int(conf.Cors.MaxAge / 1000),
	}))

	return r, nil
}
