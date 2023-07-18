package prize_handlers

import (
	"net/http"
	"xy.com/mysite/database"
	"xy.com/mysite/models/prize_models"

	"github.com/gin-gonic/gin"
)

// DrawHandler handles the draw operation.
func DrawHandler(c *gin.Context) {
	userID := c.Param("userID")

	// Fetch user's points system
	pointsSystem, err := prize_models.GetPointsSystem(database.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Perform the draw operation
	if err := pointsSystem.Draw(database.DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the points system in the database
	if err := prize_models.UpdatePointsSystem(database.DB, pointsSystem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Draw operation successful", "point": pointsSystem})
}

// ExchangeCoinsHandler handles the exchange operation.
func ExchangeCoinsHandler(c *gin.Context) {
	// Get the userID from the Gin context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Convert the userID to the desired type (e.g., string)
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get userID"})
		return
	}

	// Fetch user's points system
	pointsSystem, err := prize_models.GetPointsSystem(database.DB, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Perform the exchange operation
	if err := pointsSystem.ExchangeCoins(database.DB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the points system in the database
	if err := prize_models.UpdatePointsSystem(database.DB, pointsSystem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exchange operation successful"})
}

// GetPointsSystemHandler handles fetching the points system for a specific user.
func GetPointsSystemHandler(c *gin.Context) {
	userID := c.Param("userID")

	// Fetch user's points system
	pointsSystem, err := prize_models.GetPointsSystem(database.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pointsSystem)
}
