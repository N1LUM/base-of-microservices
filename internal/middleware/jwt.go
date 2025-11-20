package middleware

import (
	"net/http"
	"site-constructor/internal/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtManager.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}
