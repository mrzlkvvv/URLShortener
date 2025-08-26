package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/mrzlkvvv/URLShortener/internal/config"
	"github.com/mrzlkvvv/URLShortener/internal/database"
	"github.com/mrzlkvvv/URLShortener/internal/database/postgres"
	"github.com/mrzlkvvv/URLShortener/internal/logger"
	"github.com/mrzlkvvv/URLShortener/internal/server"
)

type App struct {
	config   *config.App
	logger   *zap.Logger
	database database.Database
	server   *http.Server
}

func New() *App {
	cfg := config.MustLoad()

	logger.Init(cfg.Env, cfg.Logger)
	l := zap.L()

	l.Debug("Config readed", zap.Any("config", cfg))

	db := postgres.New(cfg.Database)

	return &App{
		config:   cfg,
		logger:   l,
		database: db,
		server:   server.New(l, db, cfg.Server),
	}
}

func (a *App) Start() {
	a.logger.Info("Starting URLShortener",
		zap.String("Env", a.config.Env),
		zap.String("LogLevel", a.config.Logger.Level),
	)

	serverErr := make(chan error, 1)

	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			serverErr <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		a.logger.Fatal("Server error", zap.Error(err))
	case <-ctx.Done():
		a.logger.Info("Shutting down")
		a.shutdown()
	}
}

func (a *App) shutdown() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		a.config.Server.Timeout+1*time.Second,
	)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		a.logger.Error("Server shutdown failed", zap.Error(err))
	}

	a.database.Shutdown()

	_ = a.logger.Sync()
}
