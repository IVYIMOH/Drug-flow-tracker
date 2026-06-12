// handlers/stock.go
package handlers

import (
	"drug-flow-tracker/db"
	"drug-flow-tracker/models"
	"drug-flow-tracker/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddStock(c *gin.Context) {
	var input models.StockEntry
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate positive numbers for pricing and storage logs
	if input.Quantity <= 0 || input.UnitPrice < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity or unit price values"})
		return
	}

	// Validate supplier source
	if !input.Source.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source must be KEMSA or PRIVATE"})
		return
	}

	_, err := db.DB.ExecContext(c.Request.Context(), `
		INSERT INTO stock_entries (hospital_id, drug_name, source, quantity, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`, input.HospitalID, input.DrugName, input.Source, input.Quantity, input.UnitPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store incoming stock entry"})
		return
	}

	priceCheck := utils.CheckHighPrice(input.DrugName, input.UnitPrice)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Stock entry added successfully",
		"high_price":      priceCheck.IsHigh,
		"price_benchmark": priceCheck.IsKnown,
	})
}

func GetStock(c *gin.Context) {
	// Query calculates active stock isolated by each distinct hospital facility
	rows, err := db.DB.QueryContext(c.Request.Context(), `
		SELECT 
			s.hospital_id,
			s.drug_name,
			COALESCE(SUM(s.quantity), 0) - COALESCE(
				(SELECT SUM(d.quantity) 
				 FROM dispensations d 
				 WHERE d.drug_name = s.drug_name AND d.hospital_id = s.hospital_id), 0
			) AS current_stock
		FROM stock_entries s
		GROUP BY s.hospital_id, s.drug_name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query available database stocks"})
		return
	}
	defer rows.Close()

	results := []gin.H{}
	for rows.Next() {
		var hospitalID int
		var drug string
		var available int
		if err := rows.Scan(&hospitalID, &drug, &available); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse record entry rows"})
			return
		}
		
		results = append(results, gin.H{
			"hospital_id":     hospitalID,
			"drug_name":       drug,
			"available_stock": available,
		})
	}
	c.JSON(http.StatusOK, results)
}
