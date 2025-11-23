package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig configures CORS behavior
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int // seconds
}

// DefaultCORSConfig returns default CORS config
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}
}

// CORSMiddleware creates a CORS middleware
func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
	allowOrigins := strings.Join(config.AllowOrigins, ", ")
	allowMethods := strings.Join(config.AllowMethods, ", ")
	allowHeaders := strings.Join(config.AllowHeaders, ", ")
	exposeHeaders := strings.Join(config.ExposeHeaders, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed
		if len(config.AllowOrigins) > 0 {
			allowed := false
			for _, allowedOrigin := range config.AllowOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}
			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		c.Header("Access-Control-Allow-Methods", allowMethods)
		c.Header("Access-Control-Allow-Headers", allowHeaders)
		c.Header("Access-Control-Expose-Headers", exposeHeaders)

		if config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if config.MaxAge > 0 {
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}

		c.Next()
	}
}

