package contracts

import (
	"context"
)

// ChangePasswordCommandRequest represents change password command request
type ChangePasswordCommandRequest struct {
	UserID         string
	CurrentPassword string
	NewPassword    string
}

// ChangePasswordCommand handles password changes
type ChangePasswordCommand interface {
	Execute(ctx context.Context, req ChangePasswordCommandRequest) error
}

