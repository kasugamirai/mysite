package prize_handlers

import (
	"net/http"
	"xy.com/mysite/database"
	"xy.com/mysite/models/prize_models"

	"github.com/gin-gonic/gin"
)

// AddCodeHandler handles adding a new code to the database.
func AddCodeHandler(c *gin.Context) {
	// Parse request
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add code to database
	if err := prize_models.AddCode(database.DB, req.Code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code added successfully"})
}

// GetCodeHandler handles getting an unused code from the database.
func GetCodeHandler(c *gin.Context) {
	// Retrieve unused code
	code, err := prize_models.GetCode(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code retrieved successfully", "code": code})
}
