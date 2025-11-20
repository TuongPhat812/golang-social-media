package persistence

import (
	domain "golang-social-media/apps/chat-service/internal/domain/message"
)

// MessageMapperImpl implements MessageMapper interface
type MessageMapperImpl struct{}

var _ MessageMapper = (*MessageMapperImpl)(nil)

// NewMessageMapper creates a new MessageMapperImpl
func NewMessageMapper() MessageMapper {
	return &MessageMapperImpl{}
}

// ToModel converts a domain Message to MessageModel
func (m *MessageMapperImpl) ToModel(msg domain.Message) MessageModel {
	return MessageModel{
		ID:         msg.ID,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt,
	}
}

// ToDomain converts a MessageModel to domain Message
func (m *MessageMapperImpl) ToDomain(model MessageModel) domain.Message {
	return domain.Message{
		ID:         model.ID,
		SenderID:   model.SenderID,
		ReceiverID: model.ReceiverID,
		Content:    model.Content,
		CreatedAt:  model.CreatedAt,
	}
}

// ToDomainList converts a slice of MessageModel to domain Messages
func (m *MessageMapperImpl) ToDomainList(models []MessageModel) []domain.Message {
	messages := make([]domain.Message, len(models))
	for i, model := range models {
		messages[i] = m.ToDomain(model)
	}
	return messages
}

// ToModelList converts a slice of domain Messages to MessageModels
func (m *MessageMapperImpl) ToModelList(messages []domain.Message) []MessageModel {
	models := make([]MessageModel, len(messages))
	for i, msg := range messages {
		models[i] = m.ToModel(msg)
	}
	return models
}

