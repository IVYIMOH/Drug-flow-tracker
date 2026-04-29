package routes

import (
	"Drug-flow-tracker/handlers"

	"github.com/gin-gonic/gin"

)

func SetuoRoutes(r *gin.Engine) {
	r.POST("/stock", handlers.AddStock)
	r.POST("/dispense", handlers.DispenseDrug)
	r.GET("/stock", handlers.GetStock)
	r.GET("/insights", handlers.GetInsights)
}