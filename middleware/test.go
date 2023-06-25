package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 解析JSON请求体
		var request struct {
			Code string `json:"code"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 检查code是否为"123"
		if request.Code != "123" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
			return
		}

		// 如果code为"123"，则继续处理请求
		c.Next()
	}
}
