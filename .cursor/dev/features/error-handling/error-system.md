# Error Handling System

## Overview

The error handling system provides a centralized way to manage errors across all services with:
- **Error Codes**: Unique identifiers for each error type
- **AppError**: Structured error with code, message, details, and HTTP status
- **Error Transformer**: Pipeline to convert internal errors to user-friendly errors
- **HTTP Middleware**: Automatic error transformation for HTTP handlers
- **gRPC Interceptor**: Automatic error transformation for gRPC handlers

## Architecture

```
┌─────────────────┐
│  Domain Layer   │  Returns AppError or standard error
└────────┬────────┘
         │
┌────────▼────────┐
│ Application     │  Returns AppError or standard error
│ Layer           │
└────────┬────────┘
         │
┌────────▼────────┐
│ Infrastructure  │  Returns standard error
│ Layer           │
└────────┬────────┘
         │
┌────────▼──────────────────────────────┐
│ Error Transformer                     │
│ - Maps errors to AppError             │
│ - Adds user-friendly messages         │
│ - Includes details in dev mode        │
└────────┬──────────────────────────────┘
         │
    ┌────┴────┐
    │         │
┌───▼───┐ ┌──▼────┐
│ HTTP  │ │ gRPC  │
│ Middle│ │ Inter │
│ ware  │ │ ceptor│
└───────┘ └───────┘
```

## Error Codes

Error codes follow the pattern `ERR_XXXX` where:
- `ERR_0XXX`: Common errors
- `ERR_1XXX`: Auth service errors
- `ERR_2XXX`: Chat service errors
- `ERR_3XXX`: Notification service errors
- `ERR_4XXX`: E-commerce service errors
- `ERR_5XXX`: Database errors
- `ERR_6XXX`: External service errors
- `ERR_7XXX`: Event bus errors

### Example Error Codes

```go
CodeEmailRequired      = "ERR_1001"
CodeEmailInvalid       = "ERR_1002"
CodeEmailAlreadyExists = "ERR_1003"
CodeUserNotFound       = "ERR_1008"
```

## AppError Structure

```go
type AppError struct {
    Code       ErrorCode              // Unique error code
    Message    string                 // User-friendly message
    Details    map[string]interface{} // Additional context
    HTTPStatus int                    // HTTP status code
    Original   error                  // Original error (for logging)
}
```

## Usage

### In Domain Layer

```go
func (u User) Validate() error {
    if strings.TrimSpace(u.Email) == "" {
        return errors.NewValidationError(errors.CodeEmailRequired, nil)
    }
    if !strings.Contains(u.Email, "@") {
        return errors.NewValidationError(errors.CodeEmailInvalid, nil)
    }
    return nil
}
```

### In Application Layer

```go
func (c *registerUserCommand) Execute(ctx context.Context, req auth.RegisterRequest) (auth.RegisterResponse, error) {
    if err := userModel.Validate(); err != nil {
        return auth.RegisterResponse{}, err // Already AppError
    }
    
    if err := c.repo.Create(userModel); err != nil {
        // Transform repository errors
        if strings.Contains(err.Error(), "duplicate") {
            return auth.RegisterResponse{}, errors.NewConflictError(errors.CodeEmailAlreadyExists)
        }
        return auth.RegisterResponse{}, err
    }
    
    return auth.RegisterResponse{...}, nil
}
```

### In HTTP Handlers

```go
router.POST("/auth/register", func(c *gin.Context) {
    var req auth.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(errors.NewInvalidRequestError("Invalid request body"))
        return
    }
    
    resp, err := h.RegisterUser.Execute(c.Request.Context(), req)
    if err != nil {
        c.Error(err) // Middleware will transform
        return
    }
    c.JSON(http.StatusCreated, resp)
})
```

### In gRPC Handlers

```go
func (h *Handler) CreateMessage(ctx context.Context, req *chatv1.CreateMessageRequest) (*chatv1.CreateMessageResponse, error) {
    msg, err := h.createMessageCmd.Execute(ctx, cmdReq)
    if err != nil {
        return nil, err // Interceptor will transform
    }
    return resp, nil
}
```

## Error Transformer

The transformer automatically maps errors to AppError:

- **GORM errors**: `gorm.ErrRecordNotFound` → `CodeNotFound`
- **Database errors**: Duplicate key → `CodeConflict`
- **Domain errors**: Already AppError → Pass through
- **Validation errors**: Pattern matching on error messages
- **Default**: Internal server error

### Development Mode

In development mode, the transformer includes:
- Original error message
- Error type information
- Additional debugging details

## HTTP Error Response

```json
{
  "code": "ERR_1001",
  "message": "Email is required.",
  "details": {
    "field": "email"
  }
}
```

## gRPC Error Response

gRPC errors include the error code in the message:
```
ERR_1001: Email is required.
```

The client can parse the error code from the message.

## Error Middleware Setup

### HTTP (Gin)

```go
devMode := os.Getenv("ENV") == "development"
transformer := errors.NewTransformer(devMode)
router.Use(errors.ErrorMiddleware(transformer))
```

### gRPC

```go
devMode := os.Getenv("ENV") == "development"
transformer := errors.NewTransformer(devMode)
server := grpc.NewServer(
    grpc.UnaryInterceptor(errors.GRPCErrorInterceptor(transformer)),
)
```

## Best Practices

1. **Domain Layer**: Return AppError for business rule violations
2. **Application Layer**: Transform repository/infrastructure errors to AppError
3. **Infrastructure Layer**: Return standard errors (transformer will handle)
4. **Handlers**: Use `c.Error(err)` for HTTP, return error for gRPC
5. **Error Codes**: Use specific error codes, not generic ones
6. **Messages**: Keep messages user-friendly, avoid technical details
7. **Details**: Add context in details map for debugging

## Adding New Error Codes

1. Add error code constant in `pkg/errors/codes.go`
2. Add message in `ErrorCode.GetMessage()`
3. Update transformer if needed for pattern matching
4. Use in domain/application layer

## Example: Adding a New Error

```go
// 1. Add code
const CodeProductNameTooLong ErrorCode = "ERR_4006"

// 2. Add message
func (c ErrorCode) GetMessage() string {
    messages := map[ErrorCode]string{
        // ...
        CodeProductNameTooLong: "Product name cannot exceed 100 characters.",
    }
    // ...
}

// 3. Use in domain
func (p Product) Validate() error {
    if len(p.Name) > 100 {
        return errors.NewValidationError(errors.CodeProductNameTooLong, nil)
    }
    return nil
}
```

