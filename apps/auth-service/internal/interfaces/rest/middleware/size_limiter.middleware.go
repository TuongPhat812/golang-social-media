package middleware

import (
	"io"
	"net/http"
	"strings"

	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

// SizeLimiterConfig configures request size limiting
type SizeLimiterConfig struct {
	// MaxBodySize in bytes (default: 10MB)
	MaxBodySize int64
	// Error message when size exceeded
	ErrorMessage string
}

// DefaultSizeLimiterConfig returns default size limiter config
func DefaultSizeLimiterConfig() SizeLimiterConfig {
	return SizeLimiterConfig{
		MaxBodySize:  10 * 1024 * 1024, // 10MB
		ErrorMessage: "Request body too large",
	}
}

// SizeLimiterMiddleware creates a request size limiter middleware
func SizeLimiterMiddleware(config SizeLimiterConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Limit request body size
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxBodySize)

		// Read body to check size
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			if err.Error() == "http: request body too large" {
				logger.Component("auth.middleware.size_limiter").
					Warn().
					Str("path", c.Request.URL.Path).
					Int64("max_size", config.MaxBodySize).
					Msg("request body too large")
				c.Error(errors.NewInvalidRequestError(config.ErrorMessage))
				c.Abort()
				return
			}
			c.Error(errors.NewInvalidRequestError("Failed to read request body"))
			c.Abort()
			return
		}

		// Check size
		if int64(len(body)) > config.MaxBodySize {
			logger.Component("auth.middleware.size_limiter").
				Warn().
				Str("path", c.Request.URL.Path).
				Int64("size", int64(len(body))).
				Int64("max_size", config.MaxBodySize).
				Msg("request body exceeds limit")
			c.Error(errors.NewInvalidRequestError(config.ErrorMessage))
			c.Abort()
			return
		}

		// Restore body for handlers
		c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

		c.Next()
	}
}

