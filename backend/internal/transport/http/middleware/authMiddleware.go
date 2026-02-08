package middleware

import (
	"eshbuket/internal/service/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || !authService.ValidateSession(sessionID) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Не авторизован"})
			return
		}
		c.Next()
	}
}
