package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/mrzlkvvv/URLShortener/internal/config"
	"github.com/mrzlkvvv/URLShortener/internal/database"
)

const (
	PgErrUniqueViolation = "23505"
)

type Database struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Database) *Database {
	l := zap.L()

	connConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		l.Fatal("Parse config error", zap.Error(err))
	}

	connConfig.MaxConns = cfg.Pool.MaxConns
	connConfig.MinConns = cfg.Pool.MinConns

	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		l.Fatal("Connect to database error", zap.Error(err))
	}

	err = pool.Ping(context.Background())
	if err != nil {
		l.Error("Database ping error", zap.Error(err))
	}

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS urls(
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_urls_alias ON urls(alias);
	`)
	if err != nil {
		l.Fatal("Create tables error", zap.Error(err))
	}

	l.Info("PostgreSQL storage initialized successfully")
	return &Database{pool: pool}
}

func (s *Database) SaveURL(ctx context.Context, alias, url string) error {
	_, err := s.pool.Exec(
		ctx,
		"INSERT INTO urls(alias, url) VALUES($1, $2)",
		alias, url,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == PgErrUniqueViolation {
			return database.ErrUrlAlreadyExists
		}

		return fmt.Errorf("save URL error: %w", err)
	}

	return nil
}

func (s *Database) GetURL(ctx context.Context, alias string) (string, error) {
	var url string

	err := s.pool.QueryRow(
		ctx,
		"SELECT url FROM urls WHERE alias = $1",
		alias,
	).Scan(&url)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", database.ErrUrlIsNotExists
		}
		return "", fmt.Errorf("get URL error: %w", err)
	}

	return url, nil
}

func (s *Database) Shutdown() {
	s.pool.Close()
}
