package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kerisai/shorterms-api/config"
	"github.com/rs/zerolog/log"
)

func CreateConnPool(c config.Config) (pool *pgxpool.Pool) {
	ctx := context.Background()
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", c.DbHost, c.DbPort, c.DbName, c.DbUser, c.DbPwd, c.DbSslmode)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	defer conn.Release()

	log.Info().Msg("Established connection to database!")
	return pool
}
