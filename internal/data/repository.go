// Package etl internal/etl/export.go
package data

import (
	"fmt"
	"sync"
	"time"
)

// MetricRepository define la interfaz para el almacenamiento de métricas.
type MetricRepository interface {
	Save(metric EnrichedMetric) error
	GetMetricsByChannel(channel string, from, to time.Time, limit, offset int) ([]EnrichedMetric, error)
	GetMetricsByFunnel(utmCampaign string, from, to time.Time, limit, offset int) ([]EnrichedMetric, error)
	GetAllMetrics() ([]EnrichedMetric, error)
}

// InMemoryRepository es una implementación del Repositorio que usa un mapa en memoria.
type InMemoryRepository struct {
	mu      sync.RWMutex
	storage map[string]EnrichedMetric
}

// NewInMemoryRepository crea una nueva instancia del repositorio en memoria.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		storage: make(map[string]EnrichedMetric),
	}
}

// Save guarda una métrica en el almacén en memoria de forma segura.
func (r *InMemoryRepository) Save(metric EnrichedMetric) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := fmt.Sprintf("%s-%s-%s", metric.Date.Format("2006-01-02"), metric.CampaignID, metric.Channel)
	r.storage[key] = metric
	return nil
}

func (r *InMemoryRepository) GetMetricsByChannel(channel string, from, to time.Time, limit, offset int) ([]EnrichedMetric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []EnrichedMetric
	for _, m := range r.storage {
		if m.Channel == channel && !m.Date.Before(from) && !m.Date.After(to) {
			filtered = append(filtered, m)
		}
	}

	// Apply pagination
	start := offset
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], nil
}

// GetMetricsByFunnel with pagination
func (r *InMemoryRepository) GetMetricsByFunnel(utmCampaign string, from, to time.Time, limit, offset int) ([]EnrichedMetric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []EnrichedMetric
	for _, m := range r.storage {
		if m.UTMCampaign == utmCampaign && !m.Date.Before(from) && !m.Date.After(to) {
			filtered = append(filtered, m)
		}
	}

	// Apply pagination
	start := offset
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], nil
}

// GetAllMetrics devuelve todas las métricas almacenadas.
func (r *InMemoryRepository) GetAllMetrics() ([]EnrichedMetric, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	allMetrics := make([]EnrichedMetric, 0, len(r.storage))
	for _, metric := range r.storage {
		allMetrics = append(allMetrics, metric)
	}
	return allMetrics, nil
}
