package persistence

import (
	"time"
)

type MessageModel struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey"`
	SenderID   string    `gorm:"column:sender_id;type:text;not null"`
	ReceiverID string    `gorm:"column:receiver_id;type:text;not null"`
	Content    string    `gorm:"column:content;type:text;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	ShardID    int       `gorm:"column:shard_id;type:integer;not null;<-:false"` // Generated column, read-only (PostgreSQL calculates automatically)
}

func (MessageModel) TableName() string {
	return "messages"
}

