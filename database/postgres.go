package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Alan-00280/EventSeat-PBL-UAS.git/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool membuat connection pool ke PostgreSQL menggunakan pgxpool.
// Konfigurasi pool (max/min conns, lifetime) diambil dari config.Env.
func NewPostgresPool(env *config.Env) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(env.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("database: gagal parse postgres config: %w", err)
	}

	poolCfg.MaxConns = env.DBMaxConns
	poolCfg.MinConns = env.DBMinConns
	poolCfg.MaxConnLifetime = env.DBMaxConnLifetime
	poolCfg.MaxConnIdleTime = env.DBMaxConnIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("database: gagal membuat postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: postgres ping gagal: %w", err)
	}

	return pool, nil
}