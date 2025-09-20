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
	// Carga las variables de entorno desde el archivo .env si existe
	godotenv.Load()

	// Inicializa la configuración con valores predeterminados o de las variables de entorno
	cfg := &Config{
		Port:       getEnv("PORT", "8080"),                         // Obtiene el puerto del servidor
		AdsAPIURL:  getEnv("ADS_API_URL", ""),                      // Obtiene la URL de la API de anuncios
		CrmAPIURL:  getEnv("CRM_API_URL", ""),                      // Obtiene la URL de la API de CRM
		SinkURL:    getEnv("SINK_URL", ""),                         // Obtiene la URL del servicio SINK
		SinkSecret: getEnv("SINK_SECRET", "admira_secret_example"), // Obtiene el secreto para SINK
	}
	return cfg, nil
}

// getEnv obtiene una variable de entorno o devuelve un valor predeterminado
func getEnv(key, fallback string) string {
	// Verifica si la variable de entorno está definida
	if value, ok := os.LookupEnv(key); ok {
		return value // Devuelve el valor de la variable de entorno
	}
	return fallback // Devuelve el valor predeterminado si no está definida
}
