// internal/etl/transform.go
package etl

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/btors/admira-etl/internal/data"
)

// Transformer contiene la lógica para transformar y combinar los datos.
type Transformer struct{}

// NewTransformer crea una nueva instancia de Transformer.
func NewTransformer() *Transformer {
	return &Transformer{}
}

// CombineAndCalculateMetrics cruza los datos de Ads y CRM y calcula las métricas.
func (t *Transformer) CombineAndCalculateMetrics(adsData []data.AdPerformance, crmData []data.Opportunity) ([]data.EnrichedMetric, error) {
	// Usamos un mapa para buscar oportunidades de CRM eficientemente por su clave UTM.
	crmMap := make(map[string][]data.Opportunity)
	for _, opp := range crmData {
		// Normalizamos los UTMs para crear una clave consistente.
		key := t.createUTMKey(opp.UTMCampaign, opp.UTMSource, opp.UTMMedium)
		crmMap[key] = append(crmMap[key], opp)
	}

	var results []data.EnrichedMetric

	// Iteramos sobre cada registro de rendimiento de anuncios.
	for _, ad := range adsData {
		key := t.createUTMKey(ad.UTMCampaign, ad.UTMSource, ad.UTMMedium)

		// Buscamos las oportunidades que coincidan con la clave UTM del anuncio.
		matchingOpportunities := crmMap[key]

		// Calculamos las métricas basadas en los datos cruzados.
		leads := len(matchingOpportunities)
		closedWon := 0
		var revenue float64
		for _, opp := range matchingOpportunities {
			if opp.Stage == "closed_won" {
				closedWon++
				revenue += opp.Amount
			}
		}

		// Parseamos la fecha del anuncio.
		adDate, err := time.Parse("2006-01-02", ad.Date)
		if err != nil {
			log.Printf("WARN: could not parse date for campaign %s: %v. Skipping record.", ad.CampaignID, err)
			continue // Si la fecha es inválida, saltamos este registro.
		}

		metric := data.EnrichedMetric{
			Date:          adDate,
			Channel:       ad.Channel,
			CampaignID:    ad.CampaignID,
			UTMCampaign:   ad.UTMCampaign,
			UTMSource:     ad.UTMSource,
			UTMMedium:     ad.UTMMedium,
			Clicks:        ad.Clicks,
			Impressions:   ad.Impressions,
			Cost:          ad.Cost,
			Leads:         leads,
			Opportunities: len(matchingOpportunities), // Asumimos que 1 oportunidad = 1 lead.
			ClosedWon:     closedWon,
			Revenue:       revenue,
		}

		// Calculamos las métricas derivadas de forma segura.
		if metric.Clicks > 0 {
			metric.CPC = metric.Cost / float64(metric.Clicks)
		}
		if metric.Leads > 0 {
			metric.CPA = metric.Cost / float64(metric.Leads)
		}
		if metric.Leads > 0 {
			// Asegúrate de que la conversión a float64 se hace en ambos números ANTES de dividir.
			metric.CVRLeadToOpp = float64(metric.Opportunities) / float64(metric.Leads)
		}
		if metric.Opportunities > 0 {
			metric.CVROppToWon = float64(metric.ClosedWon) / float64(metric.Opportunities)
		}
		if metric.Cost > 0 {
			metric.ROAS = metric.Revenue / metric.Cost
		}

		results = append(results, metric)
	}

	return results, nil
}

// createUTMKey genera una clave consistente para el cruce, normalizando los datos.
func (t *Transformer) createUTMKey(campaign, source, medium string) string {
	c := strings.ToLower(strings.TrimSpace(campaign))
	s := strings.ToLower(strings.TrimSpace(source))
	m := strings.ToLower(strings.TrimSpace(medium))
	return fmt.Sprintf("%s|%s|%s", c, s, m)
}
