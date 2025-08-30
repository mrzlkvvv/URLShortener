package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// keys for access from env
const (
	DOTENV_PATH_KEY = "DOTENV_PATH"

	APP_ENV_KEY = "APP_ENV"

	SERVER_HOST_KEY        = "SERVER_HOST"
	SERVER_PORT_KEY        = "SERVER_PORT"
	SERVER_TIMEOUT_KEY     = "SERVER_TIMEOUT"
	SERVER_IDLETIMEOUT_KEY = "SERVER_IDLETIMEOUT"

	DATABASE_USER_KEY                       = "POSTGRES_USER"
	DATABASE_PASSWORD_KEY                   = "POSTGRES_PASSWORD"
	DATABASE_HOST_KEY                       = "POSTGRES_HOST"
	DATABASE_PORT_KEY                       = "POSTGRES_PORT"
	DATABASE_DBNAME_KEY                     = "POSTGRES_DB"
	DATABASE_POOL_MAXCONNS_KEY              = "DATABASE_POOL_MAXCONNS"
	DATABASE_POOL_MINIDLECONNS_KEY          = "DATABASE_POOL_MINCONNS"
	DATABASE_POOL_MINCONNS_KEY              = "DATABASE_POOL_MINIDLECONNS"
	DATABASE_POOL_MAXCONNLIFETIME_KEY       = "DATABASE_POOL_MAXCONNLIFETIME"
	DATABASE_POOL_MAXCONNLIFETIMEJITTER_KEY = "DATABASE_POOL_MAXCONNLIFETIMEJITTER"
	DATABASE_POOL_MAXCONNIDLETIME_KEY       = "DATABASE_POOL_MAXCONNIDLETIME"
	DATABASE_POOL_HEALTHCHECKPERIOD_KEY     = "DATABASE_POOL_HEALTHCHECKPERIOD"

	LOGGER_LEVEL_KEY = "LOGGER_LEVEL"
)

func init() {
	// .env file is optional
	err := godotenv.Load(os.Getenv(DOTENV_PATH_KEY))
	if err != nil {
		log.Printf("Error loading .env file: %v; continuing without it...\n", err)
	}
}

func MustLoad() *App {
	return &App{
		Env:      mustGetEnv(APP_ENV_KEY),
		Server:   mustLoadServer(),
		Database: mustLoadDatabase(),
		Logger:   mustLoadLogger(),
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		panic(fmt.Errorf("environment variable '%s' must be set", key))
	}

	return value
}

func mustLoadServer() *Server {
	timeout, err := time.ParseDuration(mustGetEnv(SERVER_TIMEOUT_KEY))
	if err != nil {
		panic(err)
	}

	idleTimeout, err := time.ParseDuration(mustGetEnv(SERVER_IDLETIMEOUT_KEY))
	if err != nil {
		panic(err)
	}

	return &Server{
		Address:     mustGetEnv(SERVER_HOST_KEY) + ":" + mustGetEnv(SERVER_PORT_KEY),
		Timeout:     timeout,
		IdleTimeout: idleTimeout,
	}
}

func mustLoadDatabase() *Database {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		mustGetEnv(DATABASE_USER_KEY),
		mustGetEnv(DATABASE_PASSWORD_KEY),
		mustGetEnv(DATABASE_HOST_KEY),
		mustGetEnv(DATABASE_PORT_KEY),
		mustGetEnv(DATABASE_DBNAME_KEY),
	)

	maxConns, err := strconv.ParseInt(mustGetEnv(DATABASE_POOL_MAXCONNS_KEY), 10, 32)
	if err != nil {
		panic(err)
	}

	minConns, err := strconv.ParseInt(mustGetEnv(DATABASE_POOL_MINCONNS_KEY), 10, 32)
	if err != nil {
		panic(err)
	}

	minIdleConns, err := strconv.ParseInt(mustGetEnv(DATABASE_POOL_MINIDLECONNS_KEY), 10, 32)
	if err != nil {
		panic(err)
	}

	maxConnLifetime, err := time.ParseDuration(mustGetEnv(DATABASE_POOL_MAXCONNLIFETIME_KEY))
	if err != nil {
		panic(err)
	}

	maxConnLifetimeJitter, err := time.ParseDuration(mustGetEnv(DATABASE_POOL_MAXCONNLIFETIMEJITTER_KEY))
	if err != nil {
		panic(err)
	}

	maxConnIdletime, err := time.ParseDuration(mustGetEnv(DATABASE_POOL_MAXCONNIDLETIME_KEY))
	if err != nil {
		panic(err)
	}

	healthCheckPeriod, err := time.ParseDuration(mustGetEnv(DATABASE_POOL_HEALTHCHECKPERIOD_KEY))
	if err != nil {
		panic(err)
	}

	return &Database{
		DSN: dsn,
		Pool: &Pool{
			MaxConns:              int32(maxConns),
			MinConns:              int32(minConns),
			MinIdleConns:          int32(minIdleConns),
			MaxConnLifetime:       maxConnLifetime,
			MaxConnLifetimeJitter: maxConnLifetimeJitter,
			MaxConnIdleTime:       maxConnIdletime,
			HealthCheckPeriod:     healthCheckPeriod,
		},
	}
}

func mustLoadLogger() *Logger {
	return &Logger{
		Level: mustGetEnv(LOGGER_LEVEL_KEY),
	}
}
