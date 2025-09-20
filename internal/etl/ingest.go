// Package etl internal/etl/ingest.go
package etl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/btors/admira-etl/internal/data"
)

// Ingestor es una estructura que maneja la ingesta de datos desde servicios externos.
type Ingestor struct {
	adsURL string       // URL del servicio de Ads
	crmURL string       // URL del servicio de CRM
	client *http.Client // Cliente HTTP para realizar solicitudes
}

// adsAPIResponse representa la estructura de la respuesta del servicio de anuncios.
type adsAPIResponse struct {
	External struct {
		Ads struct {
			Performance []data.AdPerformance `json:"performance"` // Lista de métricas de rendimiento de anuncios
		} `json:"ads"`
	} `json:"external"`
}

// crmAPIResponse representa la estructura de la respuesta del servicio CRM.
type crmAPIResponse struct {
	External struct {
		CRM struct {
			Opportunities []data.Opportunity `json:"opportunities"`
		} `json:"crm"`
	} `json:"external"`
}

// NewIngestor crea y devuelve una nueva instancia de Ingestor.
func NewIngestor(adsURL, crmURL string) *Ingestor {
	return &Ingestor{
		adsURL: adsURL,                                  // Asigna la URL del servicio de Ads
		crmURL: crmURL,                                  // Asigna la URL del servicio de CRM
		client: &http.Client{Timeout: 10 * time.Second}, // Configura un tiempo de espera de 10 segundos para las solicitudes HTTP.
	}
}

// FetchData obtiene datos de los servicios de anuncios y CRM de forma concurrente.
func (i *Ingestor) FetchData(since *time.Time) ([]data.AdPerformance, []data.Opportunity, error) {
	var wg sync.WaitGroup // WaitGroup para esperar a que ambas solicitudes terminen
	wg.Add(2)             // Añade dos tareas al WaitGroup

	var adsData []data.AdPerformance // Variable para almacenar los datos de anuncios
	var crmData []data.Opportunity   // Variable para almacenar los datos de CRM
	var adsErr, crmErr error         // Variables para capturar errores

	// Goroutine para obtener datos del servicio de anuncios.
	go func() {
		defer wg.Done()
		var adsResponse adsAPIResponse
		if err := i.fetchAndDecode(i.adsURL, &adsResponse); err != nil {
			adsErr = fmt.Errorf("failed to fetch ads data: %w", err)
			return
		}
		adsData = adsResponse.External.Ads.Performance
	}()

	// Goroutine para obtener datos del servicio CRM.
	go func() {
		defer wg.Done()
		var crmResponse crmAPIResponse
		if err := i.fetchAndDecode(i.crmURL, &crmResponse); err != nil {
			crmErr = fmt.Errorf("failed to fetch crm data: %w", err)
			return
		}
		crmData = crmResponse.External.CRM.Opportunities
	}()

	wg.Wait() // Espera a que ambas goroutines terminen

	if adsErr != nil {
		return nil, nil, adsErr
	}
	if crmErr != nil {
		return nil, nil, crmErr
	}

	// Aplica un filtro basado en la fecha proporcionada en el parámetro 'since'.
	transformer := NewTransformer()
	if since != nil {
		adsData = transformer.FilterAdsByDate(adsData, since)
		crmData = transformer.FilterCRMByDate(crmData, since)
	}

	return adsData, crmData, nil
}

// fetchAndDecode realiza una solicitud HTTP GET y decodifica la respuesta JSON.
func (i *Ingestor) fetchAndDecode(url string, target interface{}) error {
	const maxRetries = 3                     // Número máximo de reintentos para solicitudes fallidas
	const baseDelay = 500 * time.Millisecond // Retraso base para el backoff exponencial

	// Intenta realizar la solicitud hasta el número máximo de reintentos.
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Crea un contexto con tiempo de espera para la solicitud.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Crea una nueva solicitud HTTP GET.
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Realiza la solicitud HTTP.
		resp, err := i.client.Do(req)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(baseDelay * time.Duration(1<<attempt)) // Exponential backoff
				continue
			}
			return fmt.Errorf("request failed after %d attempts: %w", attempt, err)
		}

		defer resp.Body.Close() // Cierra el cuerpo de la respuesta al finalizar

		if resp.StatusCode != http.StatusOK {
			if attempt < maxRetries {
				time.Sleep(baseDelay * time.Duration(1<<attempt)) // Exponential backoff
				continue
			}
			return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		}

		// Decodifica el cuerpo de la respuesta en el destino proporcionado.
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			if attempt < maxRetries {
				time.Sleep(baseDelay * time.Duration(1<<attempt)) // Exponential backoff
				continue
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}

		return nil // Retorna nil si la solicitud y decodificación fueron exitosas
	}

	return fmt.Errorf("exceeded maximum retries for URL: %s", url)
}
