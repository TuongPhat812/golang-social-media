package contracts

import (
	"context"

	"golang-social-media/pkg/events"
)

// NotificationPublisher publishes notification-related events
type NotificationPublisher interface {
	PublishNotificationCreated(ctx context.Context, event events.NotificationCreated) error
	PublishNotificationRead(ctx context.Context, event events.NotificationRead) error
	Close() error
}

