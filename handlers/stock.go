package handlers

import (
	"Drug-flow-tracker/db"
	"Drug-flow-tracker/models"
	"Drug-flow-tracker/utils"

	"github.com/gin-gonic/gin"
)

func AddStock(c *gin.Context) {
	var input models.StockEntry

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err := db.DB.Exec(` 
	INSERT INTO stock_entries (hospital_id, drug_name, source, quantity, unit_price)
	VALUES ($1, $2, $3, $4, $5)
	`,
		input.HospitalID,
		input.DrugName,
		input.Source,
		input.Quantity,
		input.UnitPrice,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// quick insight flag
	highPrice := utils.CheckHighPrice(input.DrugName, input.UnitPrice)

	c.JSON(200, gin.H{
		"message":    "Stock entry added successfully",
		"high_price": highPrice,
	})
}

func GetStock(c *gin.Context) {
	rows, err := db.DB.Query(`
	SELECT drug_name,
	COALESCE(SUM(quantity), 0) -
	COALESCE((
		SELECT SUM(quantity)
		FROM dispensations d
		WHERE d.drug_name = s.drug_namenil
		), 0) AS current_stock
	FROM stock_entries s
	GROUP BY drug_name
	`)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var results []gin.H

	for rows.Next() {
		var drug string
		var available int
		rows.Scan(&drug, &available)
		results = append(results, gin.H{
			"drug_name":       drug,
			"available_stock": available,
		})
	}
	c.JSON(200, results)
}
