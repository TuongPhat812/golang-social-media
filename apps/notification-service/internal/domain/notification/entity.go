package notification

import (
	"time"

	"github.com/gocql/gocql"
)

type Type string

const (
	TypeWelcome     Type = "welcome"
	TypeChatMessage Type = "chat_message"
)

type Notification struct {
	ID        gocql.UUID
	UserID    string
	Type      Type
	Title     string
	Body      string
	Metadata  map[string]string
	CreatedAt time.Time
}
