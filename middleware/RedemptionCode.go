package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"xy.com/mysite/database"
	prizeModels "xy.com/mysite/models/prize_models" // 请将此处替换为实际的包路径
)

func CheckRedemptionCode() gin.HandlerFunc {
	db := database.DB
	return func(c *gin.Context) {
		// Parse the JSON body
		var json struct {
			Code string `json:"code"`
		}

		if err := c.BindJSON(&json); err != nil {
			// If there is an error parsing the JSON, return a bad request response
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}

		// Call UseRedemptionCode
		err := prizeModels.UseRedemptionCode(db, json.Code)
		if err != nil {
			if err.Error() == "code not found" {
				// If the redemption code is not found in the database or has been used, return an error response
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or used redemption code"})
			} else if err.Error() == "code already used" {
				// If the redemption code has already been used, return an error response
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Redemption code has already been used"})
			} else {
				// If there is a different error, return an internal server error response
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			return
		}

		// If the redemption code is valid and not used, set it in the context and pass the request to the next handler
		c.Set("redemptionCode", json.Code)
		c.Next()
	}
}
