package contracts

import (
	"context"
)

// ValidateTokenQueryResponse represents validate token query response
type ValidateTokenQueryResponse struct {
	Valid  bool
	UserID string
}

// ValidateTokenQuery validates a JWT token
type ValidateTokenQuery interface {
	Execute(ctx context.Context, token string) (ValidateTokenQueryResponse, error)
}

