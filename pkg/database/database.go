// Package database provides the application's PostgreSQL connection pool.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/config"
)

// Module provides Conn. fx.Provide alone is lazy - nothing consumes Conn
// yet (no repository code exists), so without the fx.Invoke below the pool
// would never actually be constructed and startup wouldn't verify
// connectivity as required.
var Module = fx.Options(
	fx.Provide(New),
	fx.Invoke(func(Conn) {}),
)

// Conn is the subset of *pgxpool.Pool used by callers.
type Conn interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Params struct {
	fx.In
	fx.Lifecycle

	Config config.Config
}

type conn struct {
	pool *pgxpool.Pool
}

func New(p Params) (Conn, error) {
	// Errors below deliberately never include the DSN - it may contain
	// credentials.
	poolConfig, err := pgxpool.ParseConfig(p.Config.GetString("database.dsn"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse database dsn: %w", err)
	}

	if n := p.Config.GetInt("database.max_open_conns"); n > 0 {
		poolConfig.MaxConns = int32(n)
	}

	if raw := p.Config.GetString("database.max_conn_lifetime"); raw != "" {
		lifetime, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("failed to parse database.max_conn_lifetime: %w", err)
		}
		poolConfig.MaxConnLifetime = lifetime
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	p.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return &conn{pool: pool}, nil
}

func (c *conn) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return c.pool.Exec(ctx, sql, args...)
}

func (c *conn) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.pool.Query(ctx, sql, args...)
}

func (c *conn) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return c.pool.QueryRow(ctx, sql, args...)
}
