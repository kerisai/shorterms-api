package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kerisai/shorterms-api/config"
	"github.com/rs/zerolog/log"
)

func RunServer(c config.Config) {
	// extra configuration
	configureCors(c)

	r := chi.NewRouter()

	// Global middlewares
	r.Use(loggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	// Default route handlers
	r.NotFound(notFound)
	r.MethodNotAllowed(methodNotAllowed)
	r.Get("/", heartbeat)

	log.Info().Msgf("Running server on port %s in %s mode", c.Port, c.Env)
	http.ListenAndServe(":"+c.Port, r)
}