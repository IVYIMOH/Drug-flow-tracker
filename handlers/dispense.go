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

	_, err := db.DB.Exec(`
	INSERT INTO dispensations (hospital_id, drug_name, quantity)
	VALUES ($1, $2, $3)
	`,
		input.HospitalID,
		input.DrugName,
		input.Quantity,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Drug dispensed successfully",
	})
}
