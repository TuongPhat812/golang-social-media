package persistence

import (
	"context"

	"golang-social-media/apps/chat-service/internal/application/messages"
	domain "golang-social-media/apps/chat-service/internal/domain/message"

	"gorm.io/gorm"
)

var _ messages.Repository = (*MessageRepository)(nil)

type MessageRepository struct {
	db     *gorm.DB
	mapper MessageMapper
}

func NewMessageRepository(db *gorm.DB, mapper MessageMapper) *MessageRepository {
	return &MessageRepository{
		db:     db,
		mapper: mapper,
	}
}

func (r *MessageRepository) Create(ctx context.Context, msg *domain.Message) error {
	model := r.mapper.ToModel(*msg)
	// Omit shard_id as it's a GENERATED column - PostgreSQL calculates it automatically
	if err := r.db.WithContext(ctx).Omit("shard_id").Create(&model).Error; err != nil {
		return err
	}
	*msg = r.mapper.ToDomain(model)
	return nil
}
