package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/kerisai/shorterms-api/config"
	"github.com/rs/cors"
)

var allowedOrigins []string
var allowedMethods []string
var allowedHeaders []string

func configureCors(c config.Config) {
	allowedOrigins = []string{c.ClientUrl}
	allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	if c.Env == "local" || c.Env == "development" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000", "http://localhost:3001")
		allowedMethods = append(allowedMethods, "HEAD", "TRACE")
	}

	allowedHeaders = []string{
		"Accept",
		"Authorization",
		"X-Forwarded-Authorization",
		"Content-Type",
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: allowedMethods,
		AllowedHeaders: allowedHeaders,
	}).Handler(next)
}

func loggerMiddleware(next http.Handler) http.Handler {
	logger := config.HttpLogger

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()

		defer func() {
			origin := r.RemoteAddr

			logger.Info().Fields(map[string]any{
				"method": r.Method,
				"version": r.Proto,
				"status": ww.Status(),
				"origin":     origin,
				"host":       r.Host,
				"path":       r.URL.Path,
				"user_agent": r.Header.Get("User-Agent"),
				"latency_ms": time.Since(start).Nanoseconds() / 1000000.0,
			}).Msg("http request")
		}()

		next.ServeHTTP(ww, r)
	})
}