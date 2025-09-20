// Package data internal/data/models.go
package data

import "time"

type AdPerformance struct {
	Date        string  `json:"date"`
	CampaignID  string  `json:"campaign_id"`
	Channel     string  `json:"channel"`
	Clicks      int     `json:"clicks"`
	Impressions int     `json:"impressions"`
	Cost        float64 `json:"cost"`
	UTMCampaign string  `json:"utm_campaign"`
	UTMSource   string  `json:"utm_source"`
	UTMMedium   string  `json:"utm_medium"`
}

type Opportunity struct {
	OpportunityID string    `json:"opportunity_id"`
	ContactEmail  string    `json:"contact_email"`
	Stage         string    `json:"stage"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	UTMCampaign   string    `json:"utm_campaign"`
	UTMSource     string    `json:"utm_source"`
	UTMMedium     string    `json:"utm_medium"`
}

// EnrichedMetric - VERSIÓN CORREGIDA Y COMPLETA
type EnrichedMetric struct {
	Date          time.Time
	Channel       string
	CampaignID    string
	UTMCampaign   string // Se mantienen para filtros internos
	UTMSource     string // Se mantienen para filtros internos
	UTMMedium     string // Se mantienen para filtros internos
	Clicks        int
	Impressions   int
	Cost          float64
	Leads         int
	Opportunities int
	ClosedWon     int
	Revenue       float64
	CPC           float64
	CPA           float64
	CVRLeadToOpp  float64 // <--- CAMBIO: CAMPO AÑADIDO
	CVROppToWon   float64
	ROAS          float64
}
