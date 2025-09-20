package data

import "time"

// AdPerformance represents the performance metrics of an advertisement.
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

// Opportunity represents a sales opportunity.
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

// EnrichedMetric represents an enriched metric combining ad performance and opportunity data.
type EnrichedMetric struct {
	Date          time.Time
	Channel       string
	CampaignID    string
	Clicks        int
	Impressions   int
	Cost          float64
	Leads         int
	Opportunities int
	ClosedWon     int
	Revenue       float64
	CPC           float64 // Cost Per Click
	CPA           float64 // Cost Per Acquisition (Lead)
	CVROppToWon   float64 // Conversion Rate from Opportunity to Closed Won
	ROAS          float64 // Return on Ad Spend
}
