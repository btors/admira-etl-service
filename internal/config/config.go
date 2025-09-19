package config

import (
	"github.com/joho/godotenv"
	"os"
)

// Config holds the configuration values
type Config struct {
	Port      string
	AdsAPIURL string
	CrmAPIURL string
}

// Load LoadConfig loads configuration from environment variables or .env file
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		AdsAPIURL: getEnv("ADS_API_URL", ""),
		CrmAPIURL: getEnv("CRM_API_URL", ""),
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
