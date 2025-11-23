package contracts

import (
	"context"
)

// RefreshTokenCommandRequest represents refresh token command request
type RefreshTokenCommandRequest struct {
	RefreshToken string
}

// RefreshTokenCommandResponse represents refresh token command response
type RefreshTokenCommandResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds
}

// RefreshTokenCommand handles token refresh
type RefreshTokenCommand interface {
	Execute(ctx context.Context, req RefreshTokenCommandRequest) (RefreshTokenCommandResponse, error)
}

