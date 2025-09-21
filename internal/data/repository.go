// Package data internal/data/export.go
package data

import (
	"fmt"
	"sync"
	"time"
)

// MetricRepository define la interfaz para el almacenamiento de métricas.
type MetricRepository interface {
	// Save guarda una métrica en el repositorio.
	Save(metric EnrichedMetric) error
	// GetMetricsByChannel obtiene métricas filtradas por canal, rango de fechas, límite y desplazamiento.
	GetMetricsByChannel(channel string, from, to time.Time, limit, offset int) ([]EnrichedMetric, error)
	// GetMetricsByFunnel obtiene métricas filtradas por campaña UTM, rango de fechas, límite y desplazamiento.
	GetMetricsByFunnel(utmCampaign string, from, to time.Time, limit, offset int) ([]EnrichedMetric, error)
	// GetAllMetrics devuelve todas las métricas almacenadas en el repositorio.
	GetAllMetrics() ([]EnrichedMetric, error)
}

// InMemoryRepository es una implementación del Repositorio que utiliza un mapa en memoria.
type InMemoryRepository struct {
	mu      sync.RWMutex
	storage map[string]EnrichedMetric
}

// NewInMemoryRepository crea una nueva instancia del repositorio en memoria.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		storage: make(map[string]EnrichedMetric), // Inicializa el mapa de almacenamiento.
	}
}

// Save guarda una métrica en el almacén en memoria de forma segura.
func (r *InMemoryRepository) Save(metric EnrichedMetric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Genera una clave única para la métrica basada en la fecha, ID de campaña y canal.
	key := fmt.Sprintf("%s-%s-%s", metric.Date.Format("2006-01-02"), metric.CampaignID, metric.Channel)
	r.storage[key] = metric
	return nil
}

// GetMetricsByChannel obtiene métricas filtradas por canal y rango de fechas, con paginación.
func (r *InMemoryRepository) GetMetricsByChannel(channel string, from, to time.Time, limit, offset int) ([]EnrichedMetric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []EnrichedMetric
	for _, m := range r.storage {
		// Filtra las métricas que coincidan con el canal y estén dentro del rango de fechas.
		if m.Channel == channel && !m.Date.Before(from) && !m.Date.After(to) {
			filtered = append(filtered, m)
		}
	}

	// Aplica la paginación.
	start := offset
	if start > len(filtered) {
		start = len(filtered) // Asegura que el inicio no exceda el tamaño de la lista.
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered) // Asegura que el final no exceda el tamaño de la lista.
	}

	return filtered[start:end], nil
}

// GetMetricsByFunnel obtiene métricas filtradas por campaña UTM y rango de fechas, con paginación.
func (r *InMemoryRepository) GetMetricsByFunnel(utmCampaign string, from, to time.Time, limit, offset int) ([]EnrichedMetric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []EnrichedMetric
	for _, m := range r.storage {
		// Filtra las métricas que coincidan con la campaña UTM y estén dentro del rango de fechas.
		if m.UTMCampaign == utmCampaign && !m.Date.Before(from) && !m.Date.After(to) {
			filtered = append(filtered, m)
		}
	}

	// Aplica la paginación.
	start := offset
	if start > len(filtered) {
		start = len(filtered) // Asegura que el inicio no exceda el tamaño de la lista.
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered) // Asegura que el final no exceda el tamaño de la lista.
	}

	return filtered[start:end], nil // Devuelve el segmento paginado de métricas.
}

// GetAllMetrics devuelve todas las métricas almacenadas en el repositorio.
func (r *InMemoryRepository) GetAllMetrics() ([]EnrichedMetric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Crea una lista para almacenar todas las métricas.
	allMetrics := make([]EnrichedMetric, 0, len(r.storage))
	for _, metric := range r.storage {
		allMetrics = append(allMetrics, metric)
	}
	return allMetrics, nil // Devuelve todas las métricas.
}
