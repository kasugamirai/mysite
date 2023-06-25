package prize_handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xy.com/mysite/database"
	"xy.com/mysite/models/prize_models"
)

func AddRedemptionCodeHandler(c *gin.Context) {
	// Parse request
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add code to database
	if err := prize_models.AddRedemptionCode(database.DB, req.Code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code added successfully"})
}
