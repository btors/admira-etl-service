package main

import (
	"github.com/btors/admira-etl/internal/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL: could not load config: %v", err)
	}

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	log.Printf("INFO: Server starting on port %s", cfg.Port)

	router.Run(":" + cfg.Port)
}
