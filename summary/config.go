package summary

import (
	"errors"

	"github.com/google/generative-ai-go/genai"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var ErrInvalidGeminiModel = errors.New("invalid gemini model")

var (
	db       *pgxpool.Pool
	gemini   *genai.Client
	genModel *genai.GenerativeModel
	logger   zerolog.Logger
)

func Configure(dbPool *pgxpool.Pool, geminiClient *genai.Client, geminiGenModel string) {
	if geminiGenModel == "" {
		log.Fatal().Err(ErrInvalidGeminiModel).Msg("Failed to configure summary dependencies")
	}

	db = dbPool
	gemini = geminiClient
	genModel = gemini.GenerativeModel(geminiGenModel)
	logger = log.With().Str("service", "summary").Logger()

	log.Info().Msg("Configured summary service dependencies!")
}
