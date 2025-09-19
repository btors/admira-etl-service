package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	router.Run(":8080")
}
