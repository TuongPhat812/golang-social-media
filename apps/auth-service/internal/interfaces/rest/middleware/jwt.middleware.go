package middleware

import (
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/pkg/errors"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware validates JWT token and extracts user ID, roles, and permissions
func JWTAuthMiddleware(jwtService *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			c.Error(errors.NewUnauthorizedError("Missing Authorization header"))
			c.Abort()
			return
		}

		// Validate token and get full claims
		claims, err := jwtService.ValidateTokenWithClaims(token)
		if err != nil {
			c.Error(errors.NewUnauthorizedError("Invalid or expired token"))
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)
		c.Set("token", token)
		c.Next()
	}
}

// extractTokenFromHeader extracts Bearer token from Authorization header
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	// Extract Bearer token
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

