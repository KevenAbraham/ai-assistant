package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/KevenAbraham/ai-assistant/internal/config"
)

// DB is a thin wrapper around a pgx connection.
// For production use, replace with pgxpool.Pool once puddle/v2 is available.
type DB struct {
	conn *pgx.Conn
}

func NewDB(ctx context.Context, cfg *config.Config) (*DB, error) {
	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx) //nolint:errcheck
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Conn() *pgx.Conn { return db.conn }

func (db *DB) Close(ctx context.Context) error { return db.conn.Close(ctx) }
