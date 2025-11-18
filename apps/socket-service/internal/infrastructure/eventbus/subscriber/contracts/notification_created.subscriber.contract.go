package contracts

import (
	"context"
)

// NotificationCreatedSubscriber subscribes to NotificationCreated events
type NotificationCreatedSubscriber interface {
	Consume(ctx context.Context)
	Close() error
}
