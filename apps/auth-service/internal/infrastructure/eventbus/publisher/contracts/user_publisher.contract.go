package contracts

import (
	"context"

	"golang-social-media/pkg/events"
)

// UserPublisher publishes user-related events
type UserPublisher interface {
	PublishUserCreated(ctx context.Context, event events.UserCreated) error
	Close() error
}

