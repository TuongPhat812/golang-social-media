package persistence

import (
	"context"

	"golang-social-media/apps/chat-service/internal/application/messages"
	domain "golang-social-media/apps/chat-service/internal/domain/message"

	"gorm.io/gorm"
)

var _ messages.Repository = (*MessageRepository)(nil)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, msg *domain.Message) error {
	model := MessageModelFromDomain(*msg)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}
	*msg = model.ToDomain()
	return nil
}
