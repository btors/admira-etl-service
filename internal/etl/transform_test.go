package etl

import (
	"testing"
	"time"

	"github.com/btors/admira-etl/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestCombineAndCalculateMetrics(t *testing.T) {
	// Crea una nueva instancia del Transformer para realizar las pruebas.
	transformer := NewTransformer()

	// Datos de prueba para Ads.
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

	// Datos de prueba para CRM.
	crmData := []data.Opportunity{
		// Oportunidad ganada que coincide con los datos de Ads
		{
			Stage:       "closed_won",
			Amount:      750.0,
			UTMCampaign: "summer_sale",
			UTMSource:   "google",
			UTMMedium:   "cpc",
		},
		// Oportunidad perdida que coincide con los datos de Ads
		{
			Stage:       "closed_lost",
			Amount:      250.0,
			UTMCampaign: "summer_sale",
			UTMSource:   "google",
			UTMMedium:   "cpc",
		},
	}

	// Ejecuta la función a probar
	results, err := transformer.CombineAndCalculateMetrics(adsData, crmData)

	// Aserciones: verificar que los resultados son los esperados
	assert.NoError(t, err)
	// Verifica que el resultado contenga exactamente un elemento.
	assert.Len(t, results, 1)

	// Verifica los valores calculados en la métrica enriquecida.
	metric := results[0]
	expectedDate, _ := time.Parse("2006-01-02", "2025-08-01")

	assert.Equal(t, expectedDate, metric.Date)         // Verifica que la fecha sea correcta.
	assert.Equal(t, 2, metric.Leads)                   // Verifica el número de leads.
	assert.Equal(t, 1, metric.ClosedWon)               // Verifica el número de oportunidades ganadas.
	assert.Equal(t, 750.0, metric.Revenue)             // Verifica los ingresos totales.
	assert.InDelta(t, 0.5, metric.CPC, 0.001)          // Verifica el costo por clic.
	assert.InDelta(t, 25.0, metric.CPA, 0.001)         // Verifica el costo por adquisición.
	assert.InDelta(t, 1.0, metric.CVRLeadToOpp, 0.001) // Verifica la tasa de conversión de lead a oportunidad.
	assert.InDelta(t, 0.5, metric.CVROppToWon, 0.001)  // Verifica la tasa de conversión de oportunidad a ganada.
	assert.InDelta(t, 15.0, metric.ROAS, 0.001)        // Verifica el retorno sobre el gasto publicitario (ROAS).
}
