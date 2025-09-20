package main

import (
	"github.com/btors/admira-etl/internal/config"
	"github.com/btors/admira-etl/internal/etl"

	"log"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	log.Println("INFO: Starting data ingestion...")
	ingestor := etl.NewIngestor(cfg.AdsAPIURL, cfg.CrmAPIURL)
	ads, crm, err := ingestor.FetchData()
	if err != nil {
		log.Fatalf("FATAL: data ingestion failed: %v", err)
	}
	// Mensaje de log corregido para mayor claridad
	log.Printf("INFO: Ingestion successful. Fetched %d ad records and %d crm records.", len(ads), len(crm))
}
