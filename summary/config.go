package summary

import (
	"github.com/google/generative-ai-go/genai"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	db     *pgxpool.Pool
	gemini *genai.Client
	logger zerolog.Logger
)

func Configure(dbPool *pgxpool.Pool, geminiClient *genai.Client) {
	db = dbPool
	gemini = geminiClient
	logger = log.With().Str("service", "summary").Logger()

	log.Info().Msg("Configured summary service dependencies!")
}
