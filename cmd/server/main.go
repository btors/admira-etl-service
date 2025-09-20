// cmd/server/main.go
package main

import (
	"log"

	// ¡Añade el import de data!
	"github.com/btors/admira-etl/internal/config"
	"github.com/btors/admira-etl/internal/data"
	"github.com/btors/admira-etl/internal/etl"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// --- Ingesta y Transformación (como antes) ---
	log.Println("INFO: Starting data ingestion...")
	ingestor := etl.NewIngestor(cfg.AdsAPIURL, cfg.CrmAPIURL)
	ads, crm, err := ingestor.FetchData()
	if err != nil {
		log.Fatalf("FATAL: data ingestion failed: %v", err)
	}
	log.Printf("INFO: Ingestion successful. Fetched %d ad records and %d crm records.", len(ads), len(crm))

	log.Println("INFO: Starting data transformation...")
	transformer := etl.NewTransformer()
	enrichedData, err := transformer.CombineAndCalculateMetrics(ads, crm)
	if err != nil {
		log.Fatalf("FATAL: data transformation failed: %v", err)
	}
	log.Printf("INFO: Transformation successful. Generated %d enriched metrics.", len(enrichedData))

	// --- Carga (¡nuevo!) ---
	log.Println("INFO: Loading data into repository...")
	repo := data.NewInMemoryRepository()
	for _, metric := range enrichedData {
		if err := repo.Save(metric); err != nil {
			log.Printf("WARN: could not save metric for campaign %s: %v", metric.CampaignID, err)
		}
	}
	// Este log no es muy útil en la implementación actual, pero sería clave
	// si estuviéramos guardando en una base de datos real.
	log.Printf("INFO: Successfully saved %d enriched metrics to the repository.", len(enrichedData))
}
