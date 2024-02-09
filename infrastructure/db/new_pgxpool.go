package db

import (
	"context"
	"sdn_list/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewPgxpool(ctx context.Context, cfg *config.DBConfig, logger *zap.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(cfg.ConnString)

	if err != nil {
		logger.Fatal("unable to parse connection string", zap.Error(err))
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Fatal("unable to create connection pool", zap.Error(err))
	}
	if err := dbpool.Ping(ctx); err != nil {
		panic(err)
	}

	return dbpool
}
