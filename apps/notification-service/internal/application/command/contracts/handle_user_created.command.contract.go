package contracts

import (
	"context"

	"golang-social-media/pkg/events"
)

// HandleUserCreatedCommand handles UserCreated events
type HandleUserCreatedCommand interface {
	Execute(ctx context.Context, event events.UserCreated) error
}

