package routes

import (
	"drug-flow-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/stock", handlers.AddStock)
	r.POST("/dispense", handlers.DispenseDrug)
	r.GET("/stock", handlers.GetStock)
	r.GET("/insights", handlers.GetInsights)
}
