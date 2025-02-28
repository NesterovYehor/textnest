package database

import (
	"context"
	"fmt"
	"time"

	"github.com/NesterovYehor/textnest/services/auth_service/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool  *pgxpool.Pool
	Close func()
}

func New(cfg *config.DBConfig, ctx context.Context) (*DB, error) {
	pgxConfig, err := pgxpool.ParseConfig(cfg.Link)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pgxConfig.MaxConns = int32(cfg.MaxOpenConns)
	pgxConfig.MinConns = int32(cfg.MaxIdleConns)
	pgxConfig.HealthCheckPeriod = time.Duration(cfg.ConnMaxLifetime) // Adjust if needed

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool, Close: pool.Close}, nil
}
