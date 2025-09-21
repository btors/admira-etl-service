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
	adsURL string
	crmURL string
	client *http.Client
}

// adsAPIResponse representa la estructura de la respuesta del servicio de anuncios.
type adsAPIResponse struct {
	External struct {
		Ads struct {
			Performance []data.AdPerformance `json:"performance"`
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
		adsURL: adsURL,
		crmURL: crmURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchData obtiene datos de los servicios de anuncios y CRM de forma concurrente.
func (i *Ingestor) FetchData(since *time.Time) ([]data.AdPerformance, []data.Opportunity, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var adsData []data.AdPerformance
	var crmData []data.Opportunity
	var adsErr, crmErr error

	go func() {
		defer wg.Done()
		var adsResponse adsAPIResponse
		if err := i.fetchAndDecode(i.adsURL, &adsResponse); err != nil {
			adsErr = fmt.Errorf("failed to fetch ads data: %w", err)
			return
		}
		adsData = adsResponse.External.Ads.Performance
	}()

	go func() {
		defer wg.Done()
		var crmResponse crmAPIResponse
		if err := i.fetchAndDecode(i.crmURL, &crmResponse); err != nil {
			crmErr = fmt.Errorf("failed to fetch crm data: %w", err)
			return
		}
		crmData = crmResponse.External.CRM.Opportunities
	}()

	wg.Wait()

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
	const maxRetries = 3
	const baseDelay = 500 * time.Millisecond

	// Intenta realizar la solicitud hasta el número máximo de reintentos.
	for attempt := 1; attempt <= maxRetries; attempt++ {

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

		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			if attempt < maxRetries {
				time.Sleep(baseDelay * time.Duration(1<<attempt)) // Exponential backoff
				continue
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}

		return nil
	}

	return fmt.Errorf("exceeded maximum retries for URL: %s", url)
}
