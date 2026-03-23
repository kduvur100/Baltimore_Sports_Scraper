package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a connection pool to PostgreSQL
func NewPool(databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	return pool, nil
}

// RunMigrations creates required tables if they don't exist
func RunMigrations(pool *pgxpool.Pool) error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
		`CREATE TABLE IF NOT EXISTS articles (
			id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			title           TEXT NOT NULL,
			url             TEXT UNIQUE NOT NULL,
			summary         TEXT,
			source          TEXT NOT NULL,
			team            TEXT NOT NULL,
			author          TEXT,
			image_url       TEXT,
			sentiment_score NUMERIC(4,3),
			published_at    TIMESTAMPTZ NOT NULL,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_team ON articles(team)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_source ON articles(source)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_published_at ON articles(published_at DESC)`,
	}

	for _, m := range migrations {
		if _, err := pool.Exec(context.Background(), m); err != nil {
			return fmt.Errorf("migration failed [%s]: %w", m[:40], err)
		}
	}
	return nil
}
