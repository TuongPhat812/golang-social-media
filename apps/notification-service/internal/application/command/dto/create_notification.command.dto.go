package dto

import (
	"time"

	"golang-social-media/apps/notification-service/internal/domain/notification"
)

type CreateNotificationCommandRequest struct {
	UserID   string
	Type     notification.Type
	Title    string
	Body     string
	Metadata map[string]string
	Time     time.Time
}
