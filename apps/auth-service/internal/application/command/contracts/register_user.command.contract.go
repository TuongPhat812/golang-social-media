package contracts

import (
	"context"

	"golang-social-media/pkg/contracts/auth"
)

// RegisterUserCommand registers a new user
type RegisterUserCommand interface {
	Execute(ctx context.Context, req auth.RegisterRequest) (auth.RegisterResponse, error)
}

