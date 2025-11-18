package contracts

import (
	"context"

	"golang-social-media/apps/notification-service/internal/application/command/dto"
	"golang-social-media/apps/notification-service/internal/domain/notification"
)

// CreateNotificationCommand creates a new notification
type CreateNotificationCommand interface {
	Execute(ctx context.Context, req dto.CreateNotificationCommandRequest) (notification.Notification, error)
	// Handle is kept for backward compatibility, delegates to Execute
	Handle(ctx context.Context, req dto.CreateNotificationCommandRequest) (notification.Notification, error)
}
