// Package config internal/config/config.go
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config contiene los valores de configuración necesarios para la aplicación
type Config struct {
	Port       string // Puerto en el que ejecutará el servidor
	AdsAPIURL  string // URL de la API de anuncios
	CrmAPIURL  string // URL de la API de CRM
	SinkURL    string // URL del servicio SINK
	SinkSecret string // Secreto para autenticar con el servicio SINK
}

// Load carga la configuración desde variables de entorno o un archivo .env
func Load() (*Config, error) {
	godotenv.Load()

	// Inicializa la configuración con valores predeterminados o de las variables de entorno
	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		AdsAPIURL:  getEnv("ADS_API_URL", ""),
		CrmAPIURL:  getEnv("CRM_API_URL", ""),
		SinkURL:    getEnv("SINK_URL", ""),
		SinkSecret: getEnv("SINK_SECRET", "admira_secret_example"),
	}
	return cfg, nil
}

// getEnv obtiene una variable de entorno o devuelve un valor predeterminado
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
