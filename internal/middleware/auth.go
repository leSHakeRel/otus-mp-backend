package middleware

import (
	"net/http"
	"strings"

	"movie-night-planner-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService *utils.JWTService
}

func NewAuthMiddleware(jwtService *utils.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// Check if it's Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"code":    "UNAUTHORIZED",
					"message": "Authorization header must be Bearer",
				},
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"code":    "UNAUTHORIZED",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := m.jwtService.ValidateToken(parts[1])
				if err == nil {
					c.Set("userID", claims.UserID)
					c.Set("email", claims.Email)
				}
			}
		}
		c.Next()
	}
}
