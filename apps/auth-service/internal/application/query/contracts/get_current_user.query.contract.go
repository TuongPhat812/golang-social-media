package contracts

import (
	"context"
)

// GetCurrentUserQueryResponse represents get current user query response
type GetCurrentUserQueryResponse struct {
	ID    string
	Email string
	Name  string
}

// GetCurrentUserQuery gets the current user from token
type GetCurrentUserQuery interface {
	Execute(ctx context.Context, userID string) (GetCurrentUserQueryResponse, error)
}

