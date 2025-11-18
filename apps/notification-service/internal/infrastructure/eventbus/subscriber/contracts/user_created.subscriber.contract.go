package contracts

import (
	"context"
)

// UserCreatedSubscriber consumes UserCreated events
type UserCreatedSubscriber interface {
	Consume(ctx context.Context)
	Close() error
}
