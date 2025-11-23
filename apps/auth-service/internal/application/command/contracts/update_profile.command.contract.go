package contracts

import (
	"context"
)

// UpdateProfileCommandRequest represents update profile command request
type UpdateProfileCommandRequest struct {
	UserID string
	Name   string
}

// UpdateProfileCommandResponse represents update profile command response
type UpdateProfileCommandResponse struct {
	ID    string
	Email string
	Name  string
}

// UpdateProfileCommand handles user profile updates
type UpdateProfileCommand interface {
	Execute(ctx context.Context, req UpdateProfileCommandRequest) (UpdateProfileCommandResponse, error)
}

