// Package etl internal/etl/export.go
package etl

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/btors/admira-etl/internal/data"
)

// Exporter maneja la exportación de métricas.
type Exporter struct {
	sinkURL    string
	sinkSecret string
	client     *http.Client
}

// NewExporter crea una nueva instancia del Exporter.
func NewExporter(sinkURL, sinkSecret string) *Exporter {
	return &Exporter{
		sinkURL:    sinkURL,
		sinkSecret: sinkSecret,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

// ExportMetrics envía las métricas a la URL de destino con una firma.
func (e *Exporter) ExportMetrics(metrics []data.EnrichedMetric) error {
	if e.sinkURL == "" {
		log.Println("WARN: SINK_URL not configured. Skipping export.")
		return nil
	}

	// 1. Marshall el cuerpo de la petición.
	payload, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	// 2. Calcula el HMAC-SHA256 del cuerpo usando el secreto.
	mac := hmac.New(sha256.New, []byte(e.sinkSecret))
	mac.Write(payload)
	signature := hex.EncodeToString(mac.Sum(nil))

	// 3. Crea la petición HTTP.
	req, err := http.NewRequest("POST", e.sinkURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 4. Añade los headers requeridos.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", signature)

	// 5. Envía la petición.
	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to sink: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sink returned non-200 status code: %d", resp.StatusCode)
	}

	log.Printf("INFO: Successfully exported %d metrics to sink.", len(metrics))
	return nil
}
