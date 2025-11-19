package mappers

import (
	domain "golang-social-media/apps/chat-service/internal/domain/message"
	"golang-social-media/apps/chat-service/internal/infrastructure/persistence"
)

// MessageMapper maps between domain Message and persistence models
type MessageMapper struct{}

// NewMessageMapper creates a new MessageMapper
func NewMessageMapper() *MessageMapper {
	return &MessageMapper{}
}

// ToModel converts a domain Message to MessageModel
func (m *MessageMapper) ToModel(msg domain.Message) persistence.MessageModel {
	return persistence.MessageModel{
		ID:         msg.ID,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt,
	}
}

// ToDomain converts a MessageModel to domain Message
func (m *MessageMapper) ToDomain(model persistence.MessageModel) domain.Message {
	return domain.Message{
		ID:         model.ID,
		SenderID:   model.SenderID,
		ReceiverID: model.ReceiverID,
		Content:    model.Content,
		CreatedAt:  model.CreatedAt,
	}
}

// ToDomainList converts a slice of MessageModel to domain Messages
func (m *MessageMapper) ToDomainList(models []persistence.MessageModel) []domain.Message {
	messages := make([]domain.Message, len(models))
	for i, model := range models {
		messages[i] = m.ToDomain(model)
	}
	return messages
}

// ToModelList converts a slice of domain Messages to MessageModels
func (m *MessageMapper) ToModelList(messages []domain.Message) []persistence.MessageModel {
	models := make([]persistence.MessageModel, len(messages))
	for i, msg := range messages {
		models[i] = m.ToModel(msg)
	}
	return models
}

