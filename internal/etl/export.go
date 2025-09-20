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

// Exporter es una estructura que maneja la exportación de métricas a un sistema externo.
type Exporter struct {
	sinkURL    string
	sinkSecret string
	client     *http.Client
}

// NewExporter crea y devuelve una nueva instancia de Exporter.
func NewExporter(sinkURL, sinkSecret string) *Exporter {
	return &Exporter{
		sinkURL:    sinkURL,
		sinkSecret: sinkSecret,
		client:     &http.Client{Timeout: 30 * time.Second}, // Configura un tiempo de espera de 30 segundos.
	}
}

// ExportMetrics envía un conjunto de métricas al sistema de destino utilizando HMAC-SHA256 para la autenticación.
func (e *Exporter) ExportMetrics(metrics []data.EnrichedMetric) error {
	// Verifica si la URL de destino está configurada.
	if e.sinkURL == "" {
		log.Println("WARN: SINK_URL not configured. Skipping export.")
		return nil
	}

	// 1. Convierte las métricas a formato JSON.
	payload, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	// 2. Calcula la firma HMAC-SHA256 del cuerpo de la solicitud.
	mac := hmac.New(sha256.New, []byte(e.sinkSecret)) // Crea el HMAC con la clave secreta.
	mac.Write(payload)                                // Escribe el payload en el HMAC.
	signature := hex.EncodeToString(mac.Sum(nil))     // Convierte la firma a una cadena hexadecimal.

	// 3. Crea una nueva solicitud HTTP POST.
	req, err := http.NewRequest("POST", e.sinkURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 4. Configura los encabezados de la solicitud.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", signature)

	// 5. Envía la solicitud al sistema de destino.
	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to sink: %w", err)
	}
	defer resp.Body.Close()

	// Verifica si el sistema de destino devolvió un código de estado 200 (éxito).
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sink returned non-200 status code: %d", resp.StatusCode)
	}

	// Registra un mensaje indicando que las métricas se exportaron correctamente.
	log.Printf("INFO: Successfully exported %d metrics to sink.", len(metrics))
	return nil
}
