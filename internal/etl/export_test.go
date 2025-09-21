package etl

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/btors/admira-etl/internal/data"
	"github.com/stretchr/testify/assert"
)

// ParseDate convierte una cadena de fecha en formato "YYYY-MM-DD" a un objeto time.Time.
func ParseDate(dateStr string) time.Time {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return parsedDate
}

func TestExporter_ExportMetrics(t *testing.T) {
	// Crea un servidor HTTP de prueba para simular la respuesta de la API de destino (Sink).
	sinkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.NotEmpty(t, r.Header.Get("X-Signature"))
		w.WriteHeader(http.StatusOK)
	}))
	defer sinkServer.Close() // Cierra el servidor al finalizar la prueba.

	exporter := NewExporter(sinkServer.URL, "test_secret")

	metrics := []data.EnrichedMetric{
		{
			Date:       ParseDate("2025-08-01"),
			CampaignID: "C-1001",
			Channel:    "google_ads",
			Clicks:     100,
			Cost:       50.0,
			Revenue:    750.0,
			ROAS:       15.0,
		},
	}

	err := exporter.ExportMetrics(metrics)

	assert.NoError(t, err)
}
