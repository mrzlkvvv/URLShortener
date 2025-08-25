package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/mrzlkvvv/URLShortener/internal/config"
	"github.com/mrzlkvvv/URLShortener/internal/http-server/handlers"
	mwLogger "github.com/mrzlkvvv/URLShortener/internal/http-server/middleware/logger"
	"github.com/mrzlkvvv/URLShortener/internal/logger"
	"github.com/mrzlkvvv/URLShortener/internal/storage/sqlite"
)

func main() {
	cfg := config.MustLoad()

	logger.Init(cfg.Env, cfg.Logger)
	l := logger.New()
	defer l.Sync()

	l.Info("Starting URLShortener", zap.String("Env", cfg.Env), zap.String("LogLevel", cfg.Logger.Level))

	s := sqlite.New(cfg.Storage)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.Logger(l))
	router.Use(middleware.URLFormat)

	router.Get("/{alias}", handlers.Redirect(s))
	router.Post("/create", handlers.Create(s))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		l.Fatal("Failed to start server", zap.Error(err))
	}
}
