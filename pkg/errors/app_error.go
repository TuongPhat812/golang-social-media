package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with code, message, and details
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	HTTPStatus int                    `json:"-"`
	Original   error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Original)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the original error
func (e *AppError) Unwrap() error {
	return e.Original
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithOriginal sets the original error
func (e *AppError) WithOriginal(err error) *AppError {
	e.Original = err
	return e
}

// NewAppError creates a new AppError
func NewAppError(code ErrorCode, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    code.GetMessage(),
		HTTPStatus: httpStatus,
		Details:    make(map[string]interface{}),
	}
}

// NewAppErrorWithMessage creates a new AppError with custom message
func NewAppErrorWithMessage(code ErrorCode, httpStatus int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Details:    make(map[string]interface{}),
	}
}

// Common error constructors
func NewInternalError(err error) *AppError {
	return NewAppError(CodeInternalError, http.StatusInternalServerError).WithOriginal(err)
}

func NewInvalidRequestError(message string) *AppError {
	return NewAppErrorWithMessage(CodeInvalidRequest, http.StatusBadRequest, message)
}

func NewNotFoundError(code ErrorCode) *AppError {
	return NewAppError(code, http.StatusNotFound)
}

func NewUnauthorizedError() *AppError {
	return NewAppError(CodeUnauthorized, http.StatusUnauthorized)
}

func NewForbiddenError() *AppError {
	return NewAppError(CodeForbidden, http.StatusForbidden)
}

func NewValidationError(code ErrorCode, details map[string]interface{}) *AppError {
	err := NewAppError(code, http.StatusBadRequest)
	if details != nil {
		err.Details = details
	}
	return err
}

func NewConflictError(code ErrorCode) *AppError {
	return NewAppError(code, http.StatusConflict)
}

func NewTooManyRequestsError(message string) *AppError {
	if message == "" {
		return NewAppError(CodeRateLimitExceeded, http.StatusTooManyRequests)
	}
	return NewAppErrorWithMessage(CodeRateLimitExceeded, http.StatusTooManyRequests, message)
}

