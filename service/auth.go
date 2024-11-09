package service

import (
	"net/http"
	"strings"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, common.Error(
				http.StatusUnauthorized,
				"헤더에서 토큰을 찾을 수 없습니다",
			))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, common.Error(
				http.StatusUnauthorized,
				"헤더에서 토큰을 찾을 수 없습니다",
			))
			c.Abort()
			return
		}

		userId, err := ValidateToken(parts[1])
		if err != nil {
			if appErr, ok := err.(*common.AppError); ok {
				c.JSON(appErr.Code.StatusCode(), common.Error(
					appErr.Code.StatusCode(),
					appErr.Message,
				))
			} else {
				c.JSON(http.StatusUnauthorized, common.Error(
					http.StatusUnauthorized,
					"토큰 검증에 실패했습니다",
				))
			}
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}
