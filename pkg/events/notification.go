package events

import "time"

type NotificationCreated struct {
	NotificationID string
	Recipient      NotificationRecipient
	Message        string
	CreatedAt      time.Time
}

type NotificationRecipient struct {
	ID string
}
