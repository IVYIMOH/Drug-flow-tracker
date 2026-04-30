package handlers

import (
	"Drug-flow-tracker/db"

	"github.com/gin-gonic/gin"
)

func GetInsights(c *gin.Context) {
	rows, _ := d.DB Query(`
	SELECT source, COUNT(*) FROM stock_entries GROUP BY source`)

	total := 0 
	private := 0 

	for rows.Next() {
		var source string
		var count int
		rows.Scan(&source, &count)

		total += count
		if source == "PRIVATE" {
			private += count
		}
	}
	privateRatio := float64(private) / float6e(total)

	alert := ""
	if privateRatio > 0.4 {
		alert = " ⚠️ High reliance on private suppliers"
	}
	c.JSON(200, gin.H{
		"private_ratio": privateRatio,
		"alert":           alert,
	})
}