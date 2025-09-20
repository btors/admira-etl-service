package etl

import (
	"testing"
	"time"

	"github.com/btors/admira-etl/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestCombineAndCalculateMetrics(t *testing.T) {
	transformer := NewTransformer()

	// Datos de prueba
	adsData := []data.AdPerformance{
		{
			Date:        "2025-08-01",
			CampaignID:  "C-1001",
			Channel:     "google_ads",
			Clicks:      100,
			Cost:        50.0,
			UTMCampaign: "summer_sale",
			UTMSource:   "google",
			UTMMedium:   "cpc",
		},
	}
	crmData := []data.Opportunity{
		// Oportunidad ganada que coincide
		{
			Stage:       "closed_won",
			Amount:      750.0,
			UTMCampaign: "summer_sale",
			UTMSource:   "google",
			UTMMedium:   "cpc",
		},
		// Oportunidad perdida que coincide
		{
			Stage:       "closed_lost",
			Amount:      250.0,
			UTMCampaign: "summer_sale",
			UTMSource:   "google",
			UTMMedium:   "cpc",
		},
	}

	// Ejecutar la función a probar
	results, err := transformer.CombineAndCalculateMetrics(adsData, crmData)

	// Aserciones: verificar que los resultados son los esperados
	assert.NoError(t, err)
	assert.Len(t, results, 1)

	metric := results[0]
	expectedDate, _ := time.Parse("2006-01-02", "2025-08-01")

	assert.Equal(t, expectedDate, metric.Date)
	assert.Equal(t, 2, metric.Leads)
	assert.Equal(t, 1, metric.ClosedWon)
	assert.Equal(t, 750.0, metric.Revenue)
	assert.InDelta(t, 0.5, metric.CPC, 0.001)
	assert.InDelta(t, 25.0, metric.CPA, 0.001)
	assert.InDelta(t, 1.0, metric.CVRLeadToOpp, 0.001) // <-- ¡CORREGIDO!
	assert.InDelta(t, 0.5, metric.CVROppToWon, 0.001)
	assert.InDelta(t, 15.0, metric.ROAS, 0.001)
}
