package etl

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngestor_FetchData(t *testing.T) {
	// Crea un servidor de prueba para simular la API de Ads.
	adsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Responde con un código de estado HTTP 200 (OK).
		w.WriteHeader(http.StatusOK)
		// Escribe una respuesta JSON simulada con datos de rendimiento de anuncios.
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
	// Asegura que el servidor de prueba se cierre al finalizar la prueba.
	defer adsServer.Close()

	// Crea un servidor de prueba para simular la API de CRM.
	crmServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Responde con un código de estado HTTP 200 (OK).
		w.WriteHeader(http.StatusOK)
		// Escribe una respuesta JSON simulada con datos de oportunidades de CRM.
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
	// Asegura que el servidor de prueba se cierre al finalizar la prueba.
	defer crmServer.Close()

	// Crea una instancia del Ingestor con las URLs de los servidores de prueba.
	ingestor := NewIngestor(adsServer.URL, crmServer.URL)

	// Llama al metodo FetchData para obtener los datos simulados de Ads y CRM.
	adsData, crmData, err := ingestor.FetchData()

	// Verifica que no se haya producido ningún error durante la obtención de datos.
	assert.NoError(t, err)
	// Verifica que se haya recibido exactamente un elemento en los datos de Ads.
	assert.Len(t, adsData, 1)

	assert.Len(t, crmData, 1)
	// Verifica que el ID de la campaña en los datos de Ads sea el esperado.
	assert.Equal(t, "C-1001", adsData[0].CampaignID)
	// Verifica que el ID de la oportunidad en los datos de CRM sea el esperado.
	assert.Equal(t, "O-9001", crmData[0].OpportunityID)
}
