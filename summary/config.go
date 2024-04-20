package summary

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	db     *pgxpool.Pool
	logger zerolog.Logger
)

func Configure(dbPool *pgxpool.Pool) {
	db = dbPool
	logger = log.With().Str("service", "summary").Logger()

	log.Info().Msg("Configured summary service dependencies!")
}
