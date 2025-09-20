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

// Ingestor handles data ingestion from external services.
type Ingestor struct {
	adsURL string
	crmURL string
	client *http.Client
}

type adsAPIResponse struct {
	External struct {
		Ads struct {
			Performance []data.AdPerformance `json:"performance"`
		} `json:"ads"`
	} `json:"external"`
}

type crmAPIResponse struct {
	External struct {
		CRM struct {
			Opportunities []data.Opportunity `json:"opportunities"`
		} `json:"crm"`
	} `json:"external"`
}

// NewIngestor creates a new Ingestor instance.
func NewIngestor(adsURL, crmURL string) *Ingestor {
	return &Ingestor{
		adsURL: adsURL,
		crmURL: crmURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchData fetches data from both Ads and CRM services concurrently.
func (i *Ingestor) FetchData() ([]data.AdPerformance, []data.Opportunity, error) {
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

	return adsData, crmData, nil
}

// fetchAndDecode performs the HTTP GET request and decodes the JSON response.
func (i *Ingestor) fetchAndDecode(url string, target interface{}) error {
	const maxRetries = 3
	const baseDelay = 500 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Create a context with timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := i.client.Do(req)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(baseDelay * time.Duration(1<<attempt)) // Exponential backoff
				continue
			}
			return fmt.Errorf("request failed after %d attempts: %w", attempt, err)
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			if attempt < maxRetries {
				time.Sleep(baseDelay * time.Duration(1<<attempt)) // Exponential backoff
				continue
			}
			return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		}

		// Decode the response body
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			if attempt < maxRetries {
				time.Sleep(baseDelay * time.Duration(1<<attempt)) // Exponential backoff
				continue
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}

		return nil // Success
	}

	return fmt.Errorf("exceeded maximum retries for URL: %s", url)
}
