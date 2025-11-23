package contracts

import (
	"context"

	"golang-social-media/pkg/events"
)

// HandleUserCreatedCommand handles UserCreated events from auth-service
type HandleUserCreatedCommand interface {
	Execute(ctx context.Context, event events.UserCreated) error
}

