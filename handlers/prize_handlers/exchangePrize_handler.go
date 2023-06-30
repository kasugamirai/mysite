package prize_handlers

import (
	"net/http"
	"xy.com/mysite/database"
	"xy.com/mysite/models/prize_models"

	"github.com/gin-gonic/gin"
)

// ExchangePrizeHandler handles the exchange of a prize_handlers.
func ExchangePrizeHandler(c *gin.Context) {
	// Parse request
	var req struct {
		UserID    string `json:"user_id"`
		PrizeName string `json:"prize_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user's points system
	pointsSystem, err := prize_models.GetPointsSystem(database.DB, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Attempt to exchange the prize_handlers
	code, err := prize_models.ExchangePrize(database.DB, req.UserID, req.PrizeName, pointsSystem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prize exchanged successfully", "code": code})
}

// CheckIfUserExchangedPrizeHandler handles a request to check if a user has exchanged a specific prize_handlers.
func CheckIfUserExchangedPrizeHandler(c *gin.Context) {
	userID := c.Param("userID")
	prizeName := c.Param("prizeName")

	hasExchanged, err := prize_models.CheckIfUserExchangedPrize(database.DB, userID, prizeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if hasExchanged {
		code, err := prize_models.GetRedemptionCode(database.DB, userID, prizeName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"hasExchanged": hasExchanged, "code": code})
	} else {
		c.JSON(http.StatusOK, gin.H{"hasExchanged": hasExchanged})
	}
}

// GetPrizeByNameHandler handles a request to get a prize_handlers by name.
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
	userID := c.Param("userID")
	prizeName := c.Param("prizeName")

	code, err := prize_models.GetRedemptionCode(database.DB, userID, prizeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redemptionCode": code})
}
