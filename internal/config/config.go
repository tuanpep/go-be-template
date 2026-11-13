package config

import (
	"fmt"
	"os"
)

// Config holds runtime settings (12-factor: env only).
type Config struct {
	HTTPAddr    string
	DatabaseURL string
	LogLevel    string
	LogFormat   string // json | text
}

// Load reads configuration from environment variables.
// DATABASE_URL has no default: a missing value must not fall back to a local credential.
func Load() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	return Config{
		HTTPAddr:    envOr("HTTP_ADDR", ":8080"),
		DatabaseURL: databaseURL,
		LogLevel:    envOr("LOG_LEVEL", "info"),
		LogFormat:   envOr("LOG_FORMAT", "json"),
	}, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
