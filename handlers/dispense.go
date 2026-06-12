// handlers/dispense.go
package handlers

import (
	"Drug-flow-tracker/db"
	"Drug-flow-tracker/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DispenseDrug(c *gin.Context) {
	var input models.Dispensation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate positive dispense quantity
	if input.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quantity must be greater than zero"})
		return
	}

	// Begin database transaction to lock rows and prevent race conditions
	tx, err := db.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize transaction"})
		return
	}
	defer tx.Rollback() // Safely rolls back changes if an error occurs

	// Calculate stock filtered strictly by the requesting hospital_id
	var available int
	err = tx.QueryRowContext(c.Request.Context(), `
		SELECT COALESCE(
			(SELECT SUM(quantity) FROM stock_entries WHERE hospital_id = $1 AND drug_name = $2), 0
		) - COALESCE(
			(SELECT SUM(quantity) FROM dispensations WHERE hospital_id = $1 AND drug_name = $2), 0
		)
	`, input.HospitalID, input.DrugName).Scan(&available)

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute stock validation"})
		return
	}

	// Enforce hard constraint check
	if input.Quantity > available {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "insufficient stock at this facility",
			"requested":       input.Quantity,
			"available_stock": available,
		})
		return
	}

	// Log transaction entry
	_, err = tx.ExecContext(c.Request.Context(), `
		INSERT INTO dispensations (hospital_id, drug_name, quantity)
		VALUES ($1, $2, $3)
	`, input.HospitalID, input.DrugName, input.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record dispensation log"})
		return
	}

	// Commit changes safely to storage engine
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist data entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drug dispensed successfully"})
}
