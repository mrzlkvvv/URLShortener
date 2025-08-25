package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/mrzlkvvv/URLShortener/internal/config"
	"github.com/mrzlkvvv/URLShortener/internal/logger"
	"github.com/mrzlkvvv/URLShortener/internal/storage"
)

var (
	stmtSaveURL *sql.Stmt
	stmtGetURL  *sql.Stmt
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Storage) *Storage {
	l := logger.New()

	db, err := sql.Open("sqlite3", cfg.DSN)
	if err != nil {
		l.Fatal("Open db error", zap.Error(err))
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);
	`)
	if err != nil {
		l.Fatal("Preparing db error", zap.Error(err))
	}

	stmtSaveURL, err = db.Prepare("INSERT INTO urls(alias, url) VALUES(?, ?)")
	if err != nil {
		l.Fatal("Preparing stmtInsertURL error", zap.Error(err))
	}

	stmtGetURL, err = db.Prepare("SELECT url FROM urls WHERE alias=?")
	if err != nil {
		l.Fatal("Preparing stmtGetURL error", zap.Error(err))
	}

	return &Storage{db: db}
}

func (s *Storage) SaveURL(alias, url string) error {
	_, err := stmtSaveURL.Exec(alias, url)

	sqliteErr, ok := err.(sqlite3.Error)
	if ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return storage.ErrUrlAlreadyExists
	}

	return err
}

func (s *Storage) GetURL(alias string) (string, error) {
	var url string

	err := stmtGetURL.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlIsNotExists
		}
		return "", fmt.Errorf("Get URL error: %s", err)
	}

	return url, nil
}
