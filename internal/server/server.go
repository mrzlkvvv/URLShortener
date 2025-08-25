package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/mrzlkvvv/URLShortener/internal/config"
	"github.com/mrzlkvvv/URLShortener/internal/server/handlers"
	mwLogger "github.com/mrzlkvvv/URLShortener/internal/server/middleware/logger"
	"github.com/mrzlkvvv/URLShortener/internal/storage"
)

func New(l *zap.Logger, s storage.Storage, cfg *config.HTTPServer) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.Logger(l))
	router.Use(middleware.URLFormat)

	router.Get("/{alias}", handlers.Redirect(s))
	router.Post("/create", handlers.Create(s))

	return &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
}
