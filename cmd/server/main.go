package main

import (
	"log"

	"github.com/btors/admira-etl/internal/api"
	"github.com/btors/admira-etl/internal/config"
	"github.com/btors/admira-etl/internal/data"
	"github.com/btors/admira-etl/internal/etl"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	// 2. Inicializar dependencias
	repo := data.NewInMemoryRepository()
	ingestor := etl.NewIngestor(cfg.AdsAPIURL, cfg.CrmAPIURL)
	transformer := etl.NewTransformer()
	exporter := etl.NewExporter(cfg.SinkURL, cfg.SinkSecret)

	// 3. Inyectar dependencias en el Handler de la API
	apiHandler := api.NewHandler(repo, ingestor, transformer, exporter)

	// 4. Configurar el router y los endpoints
	router := gin.Default()

	// Endpoint de Observabilidad
	router.GET("/healthz", apiHandler.Healthz)
	router.GET("/readyz", apiHandler.Readyz)

	// Endpoint de métricas Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Endpoint de Ingesta
	router.POST("/ingest/run", apiHandler.RunIngestion)

	// Endpoints de Métricas
	router.GET("/metrics/channel", apiHandler.GetMetricsByChannel)
	router.GET("/metrics/funnel", apiHandler.GetMetricsByFunnel)

	// Endpoint de Exportación
	router.POST("/export/run", apiHandler.RunExport)

	// 5. Iniciar el servidor
	log.Printf("INFO: Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("FATAL: could not start server: %v", err)
	}
}
