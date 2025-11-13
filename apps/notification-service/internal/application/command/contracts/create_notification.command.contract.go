package contracts

import (
	"context"

	"golang-social-media/apps/notification-service/internal/application/command/dto"
	"golang-social-media/apps/notification-service/internal/domain/notification"
)

type CreateNotificationCommand interface {
	Handle(ctx context.Context, req dto.CreateNotificationCommandRequest) (notification.Notification, error)
}
