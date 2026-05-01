// handlers/dispense.go
package handlers

import (
	"Drug-flow-tracker/db"
	"Drug-flow-tracker/models"

	"github.com/gin-gonic/gin"
)

func DispenseDrug(c *gin.Context) {
	var input models.Dispensation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Guard against dispensing more than available stock
	var available int
	err := db.DB.QueryRow(`
		SELECT COALESCE(SUM(s.quantity), 0) -
			COALESCE((SELECT SUM(d.quantity) FROM dispensations d WHERE d.drug_name = $1), 0)
		FROM stock_entries s
		WHERE s.drug_name = $1
	`, input.DrugName).Scan(&available)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if input.Quantity > available {
		c.JSON(400, gin.H{
			"error":           "insufficient stock",
			"requested":       input.Quantity,
			"available_stock": available,
		})
		return
	}

	_, err = db.DB.Exec(`
		INSERT INTO dispensations (hospital_id, drug_name, quantity)
		VALUES ($1, $2, $3)
	`, input.HospitalID, input.DrugName, input.Quantity)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Drug dispensed successfully"})
}
