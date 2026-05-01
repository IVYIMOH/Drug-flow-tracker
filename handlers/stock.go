// handlers/stock.go
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

	// Validate supplier source
	if !input.Source.IsValid() {
		c.JSON(400, gin.H{"error": "source must be KEMSA or PRIVATE"})
		return
	}

	_, err := db.DB.Exec(`
		INSERT INTO stock_entries (hospital_id, drug_name, source, quantity, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`, input.HospitalID, input.DrugName, input.Source, input.Quantity, input.UnitPrice)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	priceCheck := utils.CheckHighPrice(input.DrugName, input.UnitPrice)

	c.JSON(200, gin.H{
		"message":         "Stock entry added successfully",
		"high_price":      priceCheck.IsHigh,
		"price_benchmark": priceCheck.IsKnown,
	})
}

func GetStock(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT s.drug_name,
			COALESCE(SUM(s.quantity), 0) -
			COALESCE((
				SELECT SUM(d.quantity)
				FROM dispensations d
				WHERE d.drug_name = s.drug_name
			), 0) AS current_stock
		FROM stock_entries s
		GROUP BY s.drug_name
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []gin.H
	for rows.Next() {
		var drug string
		var available int
		if err := rows.Scan(&drug, &available); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		results = append(results, gin.H{
			"drug_name":       drug,
			"available_stock": available,
		})
	}
	c.JSON(200, results)
}
