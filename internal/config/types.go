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
	MaxConns int32
	MinConns int32
}

type Logger struct {
	Level string
}
