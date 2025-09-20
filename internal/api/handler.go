// Package api internal/api/handler.go
package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/btors/admira-etl/internal/data"
	"github.com/btors/admira-etl/internal/etl"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total de solicitudes recibidas por endpoint y método.",
		},
		[]string{"endpoint", "method"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duración de las solicitudes por endpoint.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
}

// Handler contiene las dependencias y los manejadores de la API.
type Handler struct {
	repo        data.MetricRepository
	ingestor    *etl.Ingestor
	transformer *etl.Transformer
	exporter    *etl.Exporter
}

// Middleware para medir métricas Prometheus
func prometheusMiddleware(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()

		requestsTotal.WithLabelValues(endpoint, c.Request.Method).Inc()
		requestDuration.WithLabelValues(endpoint).Observe(duration)
	}
}

// NewHandler crea una nueva instancia del Handler con sus dependencias.
func NewHandler(repo data.MetricRepository, ingestor *etl.Ingestor, transformer *etl.Transformer, exporter *etl.Exporter) *Handler {
	return &Handler{
		repo:        repo,
		ingestor:    ingestor,
		transformer: transformer,
		exporter:    exporter,
	}
}

// RunIngestion es el manejador para el endpoint POST /ingest/run
func (h *Handler) RunIngestion(c *gin.Context) {
	prometheusMiddleware("/ingest/run")(c)

	log.Println("INFO: Received request to run ingestion.")

	ads, crm, err := h.ingestor.FetchData()
	if err != nil {
		log.Printf("ERROR: Data ingestion failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ingest data"})
		return
	}

	enrichedData, err := h.transformer.CombineAndCalculateMetrics(ads, crm)
	if err != nil {
		log.Printf("ERROR: Data transformation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transform data"})
		return
	}

	for _, metric := range enrichedData {
		if err := h.repo.Save(metric); err != nil {
			log.Printf("WARN: Failed to save metric for campaign %s: %v", metric.CampaignID, err)
		}
	}

	log.Printf("INFO: Ingestion process completed successfully. Processed %d metrics.", len(enrichedData))

	c.JSON(http.StatusAccepted, gin.H{"status": "Ingestion process completed successfully."})
}

// Readyz is an endpoint to verify the readiness of the service.
func (h *Handler) Readyz(c *gin.Context) {
	// Example readiness check: Verify repository is accessible
	_, err := h.repo.GetAllMetrics()
	if err != nil {
		log.Printf("ERROR: Readiness check failed: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "repository not accessible"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

// GetMetricsByChannel es el manejador para GET /metrics/channel.
func (h *Handler) GetMetricsByChannel(c *gin.Context) {
	prometheusMiddleware("/metrics/channel")(c)

	channel := c.Query("channel")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	if channel == "" || fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters: channel, from, to"})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date format, use YYYY-MM-DD"})
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date format, use YYYY-MM-DD"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'limit' parameter"})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'offset' parameter"})
		return
	}

	metrics, err := h.repo.GetMetricsByChannel(channel, from, to, limit, offset)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve metrics by channel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetMetricsByFunnel with pagination
func (h *Handler) GetMetricsByFunnel(c *gin.Context) {
	prometheusMiddleware("/metrics/funnel")(c)

	utmCampaign := c.Query("utm_campaign")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	if utmCampaign == "" || fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters: utm_campaign, from, to"})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'from' date format, use YYYY-MM-DD"})
		return
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'to' date format, use YYYY-MM-DD"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'limit' parameter"})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'offset' parameter"})
		return
	}

	metrics, err := h.repo.GetMetricsByFunnel(utmCampaign, from, to, limit, offset)
	if err != nil {
		log.Printf("ERROR: Failed to retrieve metrics by funnel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// RunExport es el manejador para el endpoint POST /export/run.
func (h *Handler) RunExport(c *gin.Context) {
	prometheusMiddleware("/export/run")(c)

	log.Println("INFO: Received request to run export.")

	// Parse the 'date' query parameter
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameter: date"})
		return
	}

	exportDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid 'date' format, use YYYY-MM-DD"})
		return
	}

	// Retrieve all metrics and filter by the specified date
	allMetrics, err := h.repo.GetAllMetrics()
	if err != nil {
		log.Printf("ERROR: Failed to retrieve metrics from repository: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data from repository"})
		return
	}

	var filteredMetrics []data.EnrichedMetric
	for _, metric := range allMetrics {
		if metric.Date.Equal(exportDate) {
			filteredMetrics = append(filteredMetrics, metric)
		}
	}

	if len(filteredMetrics) == 0 {
		log.Printf("WARN: No metrics found for the specified date: %s", dateStr)
		c.JSON(http.StatusNoContent, gin.H{"status": "No metrics found for the specified date."})
		return
	}

	// Export the filtered metrics
	if err := h.exporter.ExportMetrics(filteredMetrics); err != nil {
		log.Printf("ERROR: Export failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export data"})
		return
	}

	log.Printf("INFO: Export process completed successfully. Exported %d metrics.", len(filteredMetrics))
	c.JSON(http.StatusAccepted, gin.H{"status": "Export process completed successfully."})
}

// Healthz es un endpoint simple para verificar que el servicio está vivo.
func (h *Handler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
