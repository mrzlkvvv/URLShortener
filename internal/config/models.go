package config

import "time"

type Config struct {
	Env        string      `yaml:"env" env-required:"true"`
	HTTPServer *HTTPServer `yaml:"http_server"`
	Storage    *Storage    `yaml:"storage" env-required:"true"`
	Logger     *Logger     `yaml:"logger" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Storage struct {
	DSN string `yaml:"dsn" env-required:"true"`
}

type Logger struct {
	Level string `yaml:"level" env-required:"true"`
}
