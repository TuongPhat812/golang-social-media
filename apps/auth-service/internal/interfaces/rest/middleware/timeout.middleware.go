package middleware

import (
	"context"
	"time"

	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware creates a middleware that sets a timeout for request processing
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		// Create channel to track completion
		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		// Wait for completion or timeout
		select {
		case <-done:
			// Request completed
		case <-ctx.Done():
			// Timeout occurred
			if ctx.Err() == context.DeadlineExceeded {
				logger.Component("auth.middleware.timeout").
					Warn().
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Dur("timeout", timeout).
					Msg("request timeout exceeded")
				c.Abort()
				c.JSON(504, gin.H{"error": "Request timeout"})
			}
		}
	}
}

