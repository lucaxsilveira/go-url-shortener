package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

// LoadConfig carrega as configurações da aplicação
func LoadConfig() *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/urlshortener"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}
}
