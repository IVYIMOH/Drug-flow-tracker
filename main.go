package main

import (
	"drug-flow-tracker/db"
	"drug-flow-tracker/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
}
