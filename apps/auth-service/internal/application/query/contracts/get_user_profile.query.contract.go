package contracts

import (
	"context"

	"golang-social-media/pkg/contracts/auth"
)

// GetUserProfileQuery retrieves user profile
type GetUserProfileQuery interface {
	Execute(ctx context.Context, userID string) (auth.ProfileResponse, error)
}

