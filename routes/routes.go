package routes

import (
	"drug-flow-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Home endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Drug Flow Tracker API",
			"status":  "running",
			"version": "1.0.0",
			"endpoints": []string{
				"GET /stock",
				"POST /stock",
				"POST /dispense",
				"GET /insights",
			},
		})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// API routes
	r.POST("/stock", handlers.AddStock)
	r.POST("/dispense", handlers.DispenseDrug)
	r.GET("/stock", handlers.GetStock)
	r.GET("/insights", handlers.GetInsights)
}