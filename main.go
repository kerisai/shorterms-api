package main

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kerisai/shorterms-api/config"
	"github.com/kerisai/shorterms-api/db"
	"github.com/kerisai/shorterms-api/http"
	"github.com/kerisai/shorterms-api/summary"
	"github.com/rs/zerolog/log"
)

func main() {
	config := config.LoadConfig()
	dbPool := db.CreateConnPool(config)

	// Configure dependencies
	http.Configure(config.ClientUrl, config.Env)
	summary.Configure(dbPool)

	r := chi.NewRouter()

	// Global middlewares
	r.Use(http.LoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(http.CorsMiddleware)

	// Default route handlers
	r.NotFound(http.NotFound)
	r.MethodNotAllowed(http.MethodNotAllowed)
	r.Get("/", http.Heartbeat)

	log.Info().Msgf("Running server on port %s in %s mode", config.Port, config.Env)
	stdhttp.ListenAndServe(":"+config.Port, r)
}
