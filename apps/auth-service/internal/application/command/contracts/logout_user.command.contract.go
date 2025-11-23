package contracts

import (
	"context"
)

// LogoutUserCommandRequest represents logout command request
type LogoutUserCommandRequest struct {
	UserID string
	Token  string
}

// LogoutUserCommand handles user logout
type LogoutUserCommand interface {
	Execute(ctx context.Context, req LogoutUserCommandRequest) error
}

