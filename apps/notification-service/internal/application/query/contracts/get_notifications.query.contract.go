package contracts

import (
	"context"

	"golang-social-media/apps/notification-service/internal/domain/notification"
)

// GetNotificationsQuery retrieves notifications for a user
type GetNotificationsQuery interface {
	Execute(ctx context.Context, userID string, limit int) ([]notification.Notification, error)
}

