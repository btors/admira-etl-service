// Package config internal/config/config.go
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration values
type Config struct {
	Port       string
	AdsAPIURL  string
	CrmAPIURL  string
	SinkURL    string // Nuevo campo
	SinkSecret string // Nuevo campo
}

// Load loads configuration from environment variables or .env file
func Load() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		AdsAPIURL:  getEnv("ADS_API_URL", ""),
		CrmAPIURL:  getEnv("CRM_API_URL", ""),
		SinkURL:    getEnv("SINK_URL", ""),
		SinkSecret: getEnv("SINK_SECRET", "admira_secret_example"),
	}
	return cfg, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
