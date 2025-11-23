package contracts

import (
	"context"
)

// UserCreatedSubscriber subscribes to user.created events
type UserCreatedSubscriber interface {
	Consume(ctx context.Context)
	Close() error
}

