package notification

import "time"

type Notification struct {
	ID          string
	RecipientID string
	Message     string
	CreatedAt   time.Time
}
