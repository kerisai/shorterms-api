package main

import (
	"context"
	stdhttp "net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kerisai/shorterms-api/config"
	"github.com/kerisai/shorterms-api/http"
	"github.com/kerisai/shorterms-api/summary"
	"github.com/rs/zerolog/log"
)

type operation func(ctx context.Context) error

func gracefulShutdown(ctx context.Context, timeout time.Duration, ops map[string]operation) <-chan struct{} {
	wait := make(chan struct{})

	go func() {
		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info().Msg("App shutdown")

		timeoutFn := time.AfterFunc(timeout, func() {
			log.Fatal().Int("timeout", int(timeout)).Msg("Shutdown timeout, force exit")
		})

		defer timeoutFn.Stop()

		var shutdownWg sync.WaitGroup

		for name, op := range ops {
			shutdownWg.Add(1)

			key, operation := name, op
			go func() {
				defer shutdownWg.Done()

				log.Info().Str("operation", key).Msg("Operation shutdown")
				if err := operation(ctx); err != nil {
					log.Err(err).Msg("Failed to shutdown operation")
					return
				}
				log.Info().Str("operation", key).Msg("Operation shutdown successfully")
			}()
		}

		shutdownWg.Wait()
		close(wait)
	}()

	return wait
}

func main() {
	cfg := config.LoadConfig()
	dbPool := config.CreateDBConnPool(cfg)
	gemini := config.CreateGeminiClient(cfg.GeminiApiKey)

	// Configure dependencies
	http.Configure(cfg.ClientUrl, cfg.Env)
	summary.Configure(dbPool, gemini, cfg.GeminiGenModel)

	r := chi.NewRouter()

	// Global middlewares
	r.Use(http.LoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(http.CorsMiddleware)

	// Default route handlers
	r.NotFound(http.NotFound)
	r.MethodNotAllowed(http.MethodNotAllowed)
	r.Get("/", http.Heartbeat)

	r.Mount("/summaries", summary.Router())

	httpServer := stdhttp.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Info().Msgf("Running server on port %s in %s mode", cfg.Port, cfg.Env)
		httpServer.ListenAndServe()
	}()

	wait := gracefulShutdown(context.Background(), 60*time.Second, map[string]operation{
		"database-shutdown": func(ctx context.Context) error {
			dbPool.Close()
			return nil
		},
		"http-server-shutdown": func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	})

	<-wait
}
