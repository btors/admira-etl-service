package etl

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngestor_FetchData(t *testing.T) {
	// Mock Ads API response
	adsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"external": {
				"ads": {
					"performance": [
						{
							"date": "2025-08-01",
							"campaign_id": "C-1001",
							"channel": "google_ads",
							"clicks": 100,
							"cost": 50.0,
							"utm_campaign": "summer_sale",
							"utm_source": "google",
							"utm_medium": "cpc"
						}
					]
				}
			}
		}`))
	}))
	defer adsServer.Close()

	// Mock CRM API response
	crmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"external": {
				"crm": {
					"opportunities": [
						{
							"opportunity_id": "O-9001",
							"stage": "closed_won",
							"amount": 750.0,
							"utm_campaign": "summer_sale",
							"utm_source": "google",
							"utm_medium": "cpc"
						}
					]
				}
			}
		}`))
	}))
	defer crmServer.Close()

	ingestor := NewIngestor(adsServer.URL, crmServer.URL)

	adsData, crmData, err := ingestor.FetchData()

	assert.NoError(t, err)
	assert.Len(t, adsData, 1)
	assert.Len(t, crmData, 1)
	assert.Equal(t, "C-1001", adsData[0].CampaignID)
	assert.Equal(t, "O-9001", crmData[0].OpportunityID)
}
