package data

import (
	"fmt"
	"sync"
)

// Repository define la interfaz para nuestras operaciones de almacenamiento.
// Usar una interfaz nos permitirá cambiar la implementación (ej. a una DB) en el futuro.
type Repository interface {
	Save(metric EnrichedMetric) error
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
	// Bloqueamos el mapa para una escritura segura (exclusiva).
	r.mu.Lock()
	// Defer asegura que el bloqueo se libere al final de la función, incluso si hay un error.
	defer r.mu.Unlock()

	// Creamos una clave única para esta métrica para manejar la idempotencia.
	// Si volvemos a procesar los mismos datos, simplemente se sobrescribirán.
	key := fmt.Sprintf("%s-%s-%s", metric.Date.Format("2006-01-02"), metric.CampaignID, metric.Channel)
	r.storage[key] = metric

	return nil
}
