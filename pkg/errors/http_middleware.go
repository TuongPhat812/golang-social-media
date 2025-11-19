package errors

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang-social-media/pkg/logger"
)

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ErrorMiddleware creates a Gin middleware that transforms errors to user-friendly responses
func ErrorMiddleware(transformer *Transformer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			appErr := transformer.Transform(err)

			// Log the error
			log := logger.Component("http.error_middleware")
			logError(log, appErr, c)

			// Return error response
			c.JSON(appErr.HTTPStatus, ErrorResponse{
				Code:    string(appErr.Code),
				Message: appErr.Message,
				Details: appErr.Details,
			})
			c.Abort()
			return
		}
	}
}

// ErrorHandler is a helper function to handle errors in handlers
func ErrorHandler(transformer *Transformer) func(*gin.Context, error) {
	return func(c *gin.Context, err error) {
		if err == nil {
			return
		}

		appErr := transformer.Transform(err)

		// Log the error
		log := logger.Component("http.error_handler")
		logError(log, appErr, c)

		// Return error response
		c.JSON(appErr.HTTPStatus, ErrorResponse{
			Code:    string(appErr.Code),
			Message: appErr.Message,
			Details: appErr.Details,
		})
		c.Abort()
	}
}

func logError(log *zerolog.Logger, appErr *AppError, c *gin.Context) {
	event := log.Error().
		Str("error_code", string(appErr.Code)).
		Str("error_message", appErr.Message).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Int("status", appErr.HTTPStatus)

	if appErr.Original != nil {
		event = event.Err(appErr.Original)
	}

	if len(appErr.Details) > 0 {
		event = event.Interface("details", appErr.Details)
	}

	event.Msg("request error")
}

