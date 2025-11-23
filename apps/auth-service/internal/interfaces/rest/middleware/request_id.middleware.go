package middleware

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"
const RequestIDHeader = "X-Request-ID"

// RequestIDMiddleware creates a middleware that generates a unique request ID
// and adds it to the context and response headers
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			// Generate new request ID
			requestID = uuid.New().String()
		}

		// Set in context
		c.Set(RequestIDKey, requestID)

		// Set in response header
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID retrieves the request ID from context
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(RequestIDKey); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}

