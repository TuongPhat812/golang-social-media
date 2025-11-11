package events

import "time"

type ChatCreated struct {
	Message   ChatMessage
	CreatedAt time.Time
}

type ChatMessage struct {
	ID         string
	SenderID   string
	ReceiverID string
	Content    string
	CreatedAt  time.Time
}
