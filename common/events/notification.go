package events

import (
	"time"

	"github.com/myself/golang-social-media/common/domain/user"
)

type NotificationCreated struct {
	NotificationID string
	Recipient      user.User
	Message        string
	CreatedAt      time.Time
}
