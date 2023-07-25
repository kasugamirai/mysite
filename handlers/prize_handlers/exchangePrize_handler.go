package prize_handlers

import (
	"net/http"
	"xy.com/mysite/database"
	"xy.com/mysite/models/prize_models"

	"github.com/gin-gonic/gin"
)

// ExchangePrizeHandler handles the exchange of a prize.
func ExchangePrizeHandler(c *gin.Context) {
	// Get the userID from the Gin context
	userID, exists := getUserID(c)
	if !exists {
		return
	}

	// Parse request
	var req struct {
		PrizeName string `json:"prize_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user's points system
	pointsSystem, err := prize_models.GetPointsSystem(database.DB, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Attempt to exchange the prize
	code, err := prize_models.ExchangePrize(database.DB, userID, req.PrizeName, pointsSystem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prize exchanged successfully", "code": code})
}

// CheckIfUserExchangedPrizeHandler handles a request to check if a user has exchanged a specific prize.
func CheckIfUserExchangedPrizeHandler(c *gin.Context) {
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

	prizeName := c.Param("prizeName")

	hasExchanged, err := prize_models.CheckIfUserExchangedPrize(database.DB, userIDStr, prizeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if hasExchanged {
		code, err := prize_models.GetRedemptionCode(database.DB, userIDStr, prizeName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"hasExchanged": hasExchanged, "code": code})
	} else {
		c.JSON(http.StatusOK, gin.H{"hasExchanged": hasExchanged})
	}
}

// GetPrizeByNameHandler handles a request to get a prize by name.
func GetPrizeByNameHandler(c *gin.Context) {
	prizeName := c.Param("prizeName")

	prize, err := prize_models.GetPrizeByName(database.DB, prizeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prize)
}

// GetRedemptionCodeHandler handles a request to get the redemption code for a specific user and prize.
func GetRedemptionCodeHandler(c *gin.Context) {
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

	prizeName := c.Param("prizeName")

	code, err := prize_models.GetRedemptionCode(database.DB, userIDStr, prizeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redemptionCode": code})
}
