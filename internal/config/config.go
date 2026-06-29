package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	InfuraAPIKey  string
	InfuraNetwork string
	Port          string
}

func Load() (Config, error) {
	// Load values from .env if present, we ignore the error if the file is missing.
	_ = godotenv.Load()

	cfg := Config{
		InfuraAPIKey:  os.Getenv("INFURA_API_KEY"),
		InfuraNetwork: os.Getenv("INFURA_NETWORK"),
		Port:          os.Getenv("PORT"),
	}

	if cfg.InfuraAPIKey == "" {
		return Config{}, fmt.Errorf("INFURA_API_KEY is empty")
	}
	if cfg.InfuraNetwork == "" {
		return Config{}, fmt.Errorf("INFURA_NETWORK is empty")
	}
	if cfg.Port == "" {
		// Default port for the demo; can be overridden via env.
		cfg.Port = "8080"
	}

	return cfg, nil
}
