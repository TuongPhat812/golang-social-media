package persistence

import (
	"time"

	domain "golang-social-media/apps/chat-service/internal/domain/message"
)

type MessageModel struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey"`
	SenderID   string    `gorm:"column:sender_id;type:text;not null"`
	ReceiverID string    `gorm:"column:receiver_id;type:text;not null"`
	Content    string    `gorm:"column:content;type:text;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
}

func (MessageModel) TableName() string {
	return "messages"
}

func MessageModelFromDomain(msg domain.Message) MessageModel {
	return MessageModel{
		ID:         msg.ID,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt,
	}
}

func (m MessageModel) ToDomain() domain.Message {
	return domain.Message{
		ID:         m.ID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
		Content:    m.Content,
		CreatedAt:  m.CreatedAt,
	}
}
