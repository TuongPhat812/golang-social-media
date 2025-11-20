package mappers

import (
	domain "golang-social-media/apps/chat-service/internal/domain/message"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
)

// MessageDTOMapper defines the contract for mapping between domain Message and gRPC DTOs
type MessageDTOMapper interface {
	FromCreateMessageRequest(req *chatv1.CreateMessageRequest) domain.Message
	ToCreateMessageResponse(msg domain.Message) *chatv1.CreateMessageResponse
	ToMessage(msg domain.Message) *chatv1.Message
	ToMessageList(messages []domain.Message) []*chatv1.Message
}


