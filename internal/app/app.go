package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/mrzlkvvv/URLShortener/internal/config"
	"github.com/mrzlkvvv/URLShortener/internal/logger"
	"github.com/mrzlkvvv/URLShortener/internal/server"
	"github.com/mrzlkvvv/URLShortener/internal/storage"
	"github.com/mrzlkvvv/URLShortener/internal/storage/sqlite"
)

type App struct {
	config  *config.Config
	logger  *zap.Logger
	storage storage.Storage
	server  *http.Server
}

func New() *App {
	cfg := config.MustLoad()
	logger.Init(cfg.Env, cfg.Logger)

	l := zap.L()
	s := sqlite.New(cfg.Storage)

	return &App{
		config:  cfg,
		logger:  l,
		storage: s,
		server:  server.New(l, s, cfg.HTTPServer),
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
		a.config.HTTPServer.Timeout+1*time.Second,
	)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		a.logger.Error("Server shutdown failed", zap.Error(err))
	}

	err = a.storage.Shutdown()
	if err != nil {
		a.logger.Error("Storage shutdown failed", zap.Error(err))
	}

	_ = a.logger.Sync()
}
