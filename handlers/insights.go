// handlers/insights.go
package handlers

import (
	"Drug-flow-tracker/db"

	"github.com/gin-gonic/gin"
)

func GetInsights(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT source, COALESCE(SUM(quantity), 0)
		FROM stock_entries
		GROUP BY source
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var total, private int
	for rows.Next() {
		var source string
		var qty int
		if err := rows.Scan(&source, &qty); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		total += qty
		if source == "PRIVATE" {
			private += qty
		}
	}

	var privateRatio float64
	if total > 0 {
		privateRatio = float64(private) / float64(total)
	}

	alert := ""
	if privateRatio > 0.4 {
		alert = "⚠️ High reliance on private suppliers"
	}

	c.JSON(200, gin.H{
		"total_units":   total,
		"private_units": private,
		"private_ratio": privateRatio,
		"alert":         alert,
	})
}
