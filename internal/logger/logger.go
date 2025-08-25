package logger

import (
	"log"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mrzlkvvv/URLShortener/internal/config"
)

var once sync.Once

func Init(env string, cfg *config.Logger) {
	once.Do(func() {
		var config zap.Config

		switch env {
		case "dev":
			config = zap.NewDevelopmentConfig()
		case "prod":
			config = zap.NewProductionConfig()
		default:
			log.Fatalf("Invalid cfg.Env value: \"%s\". Must be in {dev, prod}\n", env)
		}

		switch cfg.Level {
		case "debug":
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		case "info":
			config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		case "warn":
			config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		case "error":
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		default:
			log.Fatalf("Invalid cfg.Logger.Level value: \"%s\". Must be in {debug, info, warn, err}\n", cfg.Level)
		}

		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)

		l, err := config.Build()
		if err != nil {
			log.Fatalf("Failed to create zap.Logger: %s\n", err)
		}

		zap.ReplaceGlobals(l)
	})
}

func New() *zap.Logger {
	return zap.L()
}
