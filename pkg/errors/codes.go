package errors

// ErrorCode represents a unique error code for each error type
type ErrorCode string

const (
	// Common errors (0xxx)
	CodeInternalError     ErrorCode = "ERR_0001"
	CodeInvalidRequest    ErrorCode = "ERR_0002"
	CodeNotFound          ErrorCode = "ERR_0003"
	CodeUnauthorized      ErrorCode = "ERR_0004"
	CodeForbidden         ErrorCode = "ERR_0005"
	CodeValidationFailed  ErrorCode = "ERR_0006"
	CodeConflict          ErrorCode = "ERR_0007"
	CodeTimeout           ErrorCode = "ERR_0008"
	CodeRateLimitExceeded ErrorCode = "ERR_0009"

	// Auth service errors (1xxx)
	CodeEmailRequired      ErrorCode = "ERR_1001"
	CodeEmailInvalid       ErrorCode = "ERR_1002"
	CodeEmailAlreadyExists ErrorCode = "ERR_1003"
	CodePasswordRequired   ErrorCode = "ERR_1004"
	CodePasswordTooShort   ErrorCode = "ERR_1005"
	CodePasswordInvalid    ErrorCode = "ERR_1006"
	CodeNameRequired       ErrorCode = "ERR_1007"
	CodeUserNotFound       ErrorCode = "ERR_1008"
	CodeInvalidCredentials ErrorCode = "ERR_1009"
	CodeTokenInvalid       ErrorCode = "ERR_1010"
	CodeTokenExpired       ErrorCode = "ERR_1011"

	// Chat service errors (2xxx)
	CodeMessageContentRequired ErrorCode = "ERR_2001"
	CodeMessageContentTooLong  ErrorCode = "ERR_2002"
	CodeSenderIDRequired       ErrorCode = "ERR_2003"
	CodeReceiverIDRequired     ErrorCode = "ERR_2004"
	CodeMessageNotFound        ErrorCode = "ERR_2005"
	CodeChatNotFound           ErrorCode = "ERR_2006"

	// Notification service errors (3xxx)
	CodeNotificationNotFound ErrorCode = "ERR_3001"
	CodeNotificationFailed   ErrorCode = "ERR_3002"

	// E-commerce service errors (4xxx)
	CodeProductNotFound          ErrorCode = "ERR_4001"
	CodeProductNameRequired      ErrorCode = "ERR_4002"
	CodeProductPriceInvalid      ErrorCode = "ERR_4003"
	CodeProductStockInvalid      ErrorCode = "ERR_4004"
	CodeProductOutOfStock        ErrorCode = "ERR_4005"
	CodeOrderNotFound            ErrorCode = "ERR_4006"
	CodeOrderItemRequired        ErrorCode = "ERR_4007"
	CodeOrderItemQuantityInvalid ErrorCode = "ERR_4008"
	CodeOrderAlreadyConfirmed    ErrorCode = "ERR_4009"
	CodeOrderAlreadyCancelled    ErrorCode = "ERR_4010"
	CodeOrderCannotBeCancelled   ErrorCode = "ERR_4011"

	// Database errors (5xxx)
	CodeDatabaseConnection  ErrorCode = "ERR_5001"
	CodeDatabaseQuery       ErrorCode = "ERR_5002"
	CodeDatabaseTransaction ErrorCode = "ERR_5003"
	CodeDatabaseConstraint  ErrorCode = "ERR_5004"

	// External service errors (6xxx)
	CodeExternalServiceUnavailable ErrorCode = "ERR_6001"
	CodeExternalServiceTimeout     ErrorCode = "ERR_6002"
	CodeExternalServiceError       ErrorCode = "ERR_6003"

	// Event bus errors (7xxx)
	CodeEventPublishFailed   ErrorCode = "ERR_7001"
	CodeEventSubscribeFailed ErrorCode = "ERR_7002"
)

// GetMessage returns a user-friendly message for the error code
func (c ErrorCode) GetMessage() string {
	messages := map[ErrorCode]string{
		// Common
		CodeInternalError:     "An internal error occurred. Please try again later.",
		CodeInvalidRequest:    "The request is invalid. Please check your input.",
		CodeNotFound:          "The requested resource was not found.",
		CodeUnauthorized:      "You are not authorized to perform this action.",
		CodeForbidden:         "Access to this resource is forbidden.",
		CodeValidationFailed:  "Validation failed. Please check your input.",
		CodeConflict:          "The operation conflicts with the current state.",
		CodeTimeout:           "The request timed out. Please try again.",
		CodeRateLimitExceeded: "Rate limit exceeded. Please try again later.",

		// Auth
		CodeEmailRequired:      "Email is required.",
		CodeEmailInvalid:       "Email format is invalid.",
		CodeEmailAlreadyExists: "An account with this email already exists.",
		CodePasswordRequired:   "Password is required.",
		CodePasswordTooShort:   "Password must be at least 6 characters long.",
		CodePasswordInvalid:    "Password is incorrect.",
		CodeNameRequired:       "Name is required.",
		CodeUserNotFound:       "User not found.",
		CodeInvalidCredentials: "Invalid email or password.",
		CodeTokenInvalid:       "Invalid authentication token.",
		CodeTokenExpired:       "Authentication token has expired.",

		// Chat
		CodeMessageContentRequired: "Message content is required.",
		CodeMessageContentTooLong:  "Message content is too long.",
		CodeSenderIDRequired:       "Sender ID is required.",
		CodeReceiverIDRequired:     "Receiver ID is required.",
		CodeMessageNotFound:        "Message not found.",
		CodeChatNotFound:           "Chat not found.",

		// Notification
		CodeNotificationNotFound: "Notification not found.",
		CodeNotificationFailed:   "Failed to create notification.",

		// E-commerce
		CodeProductNotFound:          "Product not found.",
		CodeProductNameRequired:      "Product name is required.",
		CodeProductPriceInvalid:      "Product price is invalid.",
		CodeProductStockInvalid:      "Product stock is invalid.",
		CodeProductOutOfStock:        "Product is out of stock.",
		CodeOrderNotFound:            "Order not found.",
		CodeOrderItemRequired:        "Order must have at least one item.",
		CodeOrderItemQuantityInvalid: "Order item quantity is invalid.",
		CodeOrderAlreadyConfirmed:    "Order has already been confirmed.",
		CodeOrderAlreadyCancelled:    "Order has already been cancelled.",
		CodeOrderCannotBeCancelled:   "Order cannot be cancelled in its current state.",

		// Database
		CodeDatabaseConnection:  "Database connection failed.",
		CodeDatabaseQuery:       "Database query failed.",
		CodeDatabaseTransaction: "Database transaction failed.",
		CodeDatabaseConstraint:  "Database constraint violation.",

		// External service
		CodeExternalServiceUnavailable: "External service is unavailable.",
		CodeExternalServiceTimeout:     "External service request timed out.",
		CodeExternalServiceError:       "External service returned an error.",

		// Event bus
		CodeEventPublishFailed:   "Failed to publish event.",
		CodeEventSubscribeFailed: "Failed to subscribe to event.",
	}

	if msg, ok := messages[c]; ok {
		return msg
	}
	return "An unknown error occurred."
}
