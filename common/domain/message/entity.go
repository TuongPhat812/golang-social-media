package message

import "time"

type Message struct {
	ID         string
	SenderID   string
	ReceiverID string
	Content    string
	CreatedAt  time.Time
}
