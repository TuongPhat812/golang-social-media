package middleware

import (
	"time"

	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

// RequestLogMiddleware creates a structured request logging middleware
func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := GetClientIP(c)
		requestID := ""
		if id, exists := c.Get(RequestIDKey); exists {
			if str, ok := id.(string); ok {
				requestID = str
			}
		}

		// Process request
		c.Next()

		// Calculate duration
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		// Build log entry
		logEvent := logger.Component("auth.middleware.request_log").
			Info().
			Str("method", method).
			Str("path", path).
			Str("query", raw).
			Str("ip", clientIP).
			Int("status", statusCode).
			Dur("latency", latency).
			Int("size", bodySize).
			Str("user_agent", c.Request.UserAgent())

		if requestID != "" {
			logEvent = logEvent.Str("request_id", requestID)
		}

		// Add user ID if authenticated
		if userID, exists := c.Get("user_id"); exists {
			logEvent = logEvent.Str("user_id", userID.(string))
		}

		// Log error if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logEvent = logEvent.Err(err).Str("error_type", err.Type.String())
			}
		}

		// Log based on status code
		if statusCode >= 500 {
			logEvent.Msg("server error")
		} else if statusCode >= 400 {
			logEvent.Msg("client error")
		} else {
			logEvent.Msg("request completed")
		}
	}
}

