package persistence

import "golang-social-media/apps/chat-service/internal/domain/message"

// MessageMapper defines the contract for mapping between domain Message and persistence models
type MessageMapper interface {
	ToModel(msg message.Message) MessageModel
	ToDomain(model MessageModel) message.Message
	ToDomainList(models []MessageModel) []message.Message
	ToModelList(messages []message.Message) []MessageModel
}


