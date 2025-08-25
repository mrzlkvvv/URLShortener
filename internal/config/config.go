package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("Env var CONFIG_PATH is not set")
	}

	_, err := os.Stat(configPath)
	if err != nil {
		log.Fatalf("Stat CONFIG_PATH error: %s", err)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Read config error: %s", err)
	}

	return &cfg
}
