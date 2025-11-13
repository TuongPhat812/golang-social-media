package events

import "time"

type Notification struct {
	ID        string            `json:"id"`
	UserID    string            `json:"userId"`
	Type      string            `json:"type"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
}

type NotificationCreated struct {
	Notification Notification `json:"notification"`
}
