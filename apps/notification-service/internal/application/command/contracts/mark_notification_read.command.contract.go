package contracts

import (
	"context"
)

// MarkNotificationReadCommand marks a notification as read
type MarkNotificationReadCommand interface {
	Execute(ctx context.Context, userID string, notificationID string) error
}

