package middleware

import (
	"net/http"
	"strings"

	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "user_id"
)

// AuthClient interface for validating tokens
type AuthClient interface {
	ValidateToken(ctx gin.Context, token string) (userID string, valid bool, err error)
}

// JWTAuthMiddleware validates JWT token via auth service gRPC and extracts user ID
func JWTAuthMiddleware(authClient AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Component("gateway.middleware.jwt").
				Warn().
				Msg("missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Component("gateway.middleware.jwt").
				Warn().
				Msg("invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token via auth service
		userID, valid, err := authClient.ValidateToken(c, tokenString)
		if err != nil {
			logger.Component("gateway.middleware.jwt").
				Error().
				Err(err).
				Msg("failed to validate token with auth service")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication service unavailable"})
			c.Abort()
			return
		}

		if !valid {
			logger.Component("gateway.middleware.jwt").
				Warn().
				Msg("invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		if userID == "" {
			logger.Component("gateway.middleware.jwt").
				Warn().
				Msg("missing user_id in token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user_id in token"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set(UserIDKey, userID)

		logger.Component("gateway.middleware.jwt").
			Debug().
			Str("user_id", userID).
			Msg("jwt authentication successful")

		c.Next()
	}
}

