package prize_handlers

import (
	"net/http"
	"xy.com/mysite/database"
	"xy.com/mysite/models/prize_models"

	"github.com/gin-gonic/gin"
)

// AddPrizeHandler handles the addition of a new prize_handlers.
func AddPrizeHandler(c *gin.Context) {
	var req struct {
		PrizeName string `json:"prize_name"`
		Cost      int    `json:"cost"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := prize_models.AddPrize(database.DB, req.PrizeName, req.Cost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Prize added successfully"})
}
