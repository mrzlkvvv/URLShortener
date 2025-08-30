package config

import "time"

type App struct {
	Env      string
	Server   *Server
	Database *Database
	Logger   *Logger
}

type Server struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type Database struct {
	DSN  string
	Pool *Pool
}

type Pool struct {
	MaxConns              int32
	MinConns              int32
	MinIdleConns          int32
	MaxConnLifetime       time.Duration
	MaxConnLifetimeJitter time.Duration
	MaxConnIdleTime       time.Duration
	HealthCheckPeriod     time.Duration
}

type Logger struct {
	Level string
}
