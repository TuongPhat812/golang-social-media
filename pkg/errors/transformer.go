package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Transformer transforms internal errors into user-friendly AppErrors
type Transformer struct {
	// Enable detailed error messages in development
	DevelopmentMode bool
}

// NewTransformer creates a new error transformer
func NewTransformer(developmentMode bool) *Transformer {
	return &Transformer{
		DevelopmentMode: developmentMode,
	}
}

// Transform converts an error into an AppError
func (t *Transformer) Transform(err error) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, return it
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Check for wrapped AppError
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// Transform based on error type
	errStr := strings.ToLower(err.Error())

	// Database errors
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
		return NewNotFoundError(CodeNotFound)
	}

	// GORM errors
	if strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "unique constraint") {
		return NewConflictError(CodeConflict).WithDetails("reason", "duplicate entry")
	}

	if strings.Contains(errStr, "foreign key constraint") {
		return NewValidationError(CodeDatabaseConstraint, map[string]interface{}{
			"reason": "referenced resource does not exist",
		})
	}

	// Domain validation errors - try to map to specific error codes
	if strings.Contains(errStr, "email") {
		if strings.Contains(errStr, "empty") || strings.Contains(errStr, "required") {
			return NewValidationError(CodeEmailRequired, nil)
		}
		if strings.Contains(errStr, "invalid") || strings.Contains(errStr, "valid") {
			return NewValidationError(CodeEmailInvalid, nil)
		}
		if strings.Contains(errStr, "already exists") || strings.Contains(errStr, "duplicate") {
			return NewConflictError(CodeEmailAlreadyExists)
		}
	}

	if strings.Contains(errStr, "password") {
		if strings.Contains(errStr, "empty") || strings.Contains(errStr, "required") {
			return NewValidationError(CodePasswordRequired, nil)
		}
		if strings.Contains(errStr, "too short") || strings.Contains(errStr, "at least") {
			return NewValidationError(CodePasswordTooShort, nil)
		}
		if strings.Contains(errStr, "incorrect") || strings.Contains(errStr, "invalid") {
			return NewValidationError(CodePasswordInvalid, nil)
		}
	}

	if strings.Contains(errStr, "name") && (strings.Contains(errStr, "empty") || strings.Contains(errStr, "required")) {
		return NewValidationError(CodeNameRequired, nil)
	}

	if strings.Contains(errStr, "message") || strings.Contains(errStr, "content") {
		if strings.Contains(errStr, "empty") || strings.Contains(errStr, "required") {
			return NewValidationError(CodeMessageContentRequired, nil)
		}
		if strings.Contains(errStr, "too long") {
			return NewValidationError(CodeMessageContentTooLong, nil)
		}
	}

	if strings.Contains(errStr, "not found") {
		if strings.Contains(errStr, "user") {
			return NewNotFoundError(CodeUserNotFound)
		}
		if strings.Contains(errStr, "message") {
			return NewNotFoundError(CodeMessageNotFound)
		}
		if strings.Contains(errStr, "product") {
			return NewNotFoundError(CodeProductNotFound)
		}
		if strings.Contains(errStr, "order") {
			return NewNotFoundError(CodeOrderNotFound)
		}
		return NewNotFoundError(CodeNotFound)
	}

	if strings.Contains(errStr, "out of stock") || strings.Contains(errStr, "insufficient stock") {
		return NewValidationError(CodeProductOutOfStock, nil)
	}

	// Context timeout
	if strings.Contains(errStr, "context deadline exceeded") || strings.Contains(errStr, "timeout") {
		return NewAppError(CodeTimeout, 408).WithOriginal(err)
	}

	// Default: internal server error
	appErr := NewInternalError(err)
	if t.DevelopmentMode {
		if appErr.Details == nil {
			appErr.Details = make(map[string]interface{})
		}
		appErr.Details["original_error"] = err.Error()
		appErr.Details["error_type"] = fmt.Sprintf("%T", err)
	}
	return appErr
}

// TransformWithCode transforms an error with a specific error code
func (t *Transformer) TransformWithCode(err error, code ErrorCode, httpStatus int) *AppError {
	if err == nil {
		return nil
	}

	// If it's already an AppError, return it
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	appErr := NewAppError(code, httpStatus).WithOriginal(err)
	if t.DevelopmentMode {
		appErr.Details["original_error"] = err.Error()
	}
	return appErr
}

