// handlers/insights.go
package handlers

import (
	"Drug-flow-tracker/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInsights(c *gin.Context) {
	// Query calculates active stock remaining in the ecosystem per source type
	rows, err := db.DB.QueryContext(c.Request.Context(), `
		SELECT 
			s.source,
			COALESCE(SUM(s.quantity), 0) - COALESCE(
				(SELECT SUM(d.quantity) 
				 FROM dispensations d 
				 JOIN stock_entries se ON d.drug_name = se.drug_name AND d.hospital_id = se.hospital_id
				 WHERE se.source = s.source), 0
			) AS active_quantity
		FROM stock_entries s
		GROUP BY s.source
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate stock insights"})
		return
	}
	defer rows.Close()

	var totalActive, privateActive int
	for rows.Next() {
		var source string
		var qty int
		if err := rows.Scan(&source, &qty); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process row metrics"})
			return
		}
		
		// Prevent negative allocations from corrupted historical datasets
		if qty < 0 {
			qty = 0
		}

		totalActive += qty
		if source == "PRIVATE" {
			privateActive += qty
		}
	}

	// Calculate dependency metric on current active holdings
	var privateRatio float64
	if totalActive > 0 {
		privateRatio = float64(privateActive) / float64(totalActive)
	}

	alert := ""
	if privateRatio > 0.4 {
		alert = "⚠️ High reliance on private suppliers. Consider optimization via KEMSA procurement channels."
	}

	c.JSON(http.StatusOK, gin.H{
		"active_total_units":   totalActive,
		"active_private_units": privateActive,
		"private_ratio":        privateRatio,
		"alert":                alert,
	})
}
