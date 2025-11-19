package mappers

import (
	domain "golang-social-media/apps/chat-service/internal/domain/message"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MessageDTOMapper maps between domain Message and gRPC DTOs
type MessageDTOMapper struct{}

// NewMessageDTOMapper creates a new MessageDTOMapper
func NewMessageDTOMapper() *MessageDTOMapper {
	return &MessageDTOMapper{}
}

// ToCreateMessageRequest converts gRPC CreateMessageRequest to domain Message
// Note: ID and CreatedAt will be set by application layer
func (m *MessageDTOMapper) FromCreateMessageRequest(req *chatv1.CreateMessageRequest) domain.Message {
	return domain.Message{
		SenderID:   req.GetSenderId(),
		ReceiverID: req.GetReceiverId(),
		Content:    req.GetContent(),
		// ID and CreatedAt will be set by application layer
	}
}

// ToCreateMessageResponse converts domain Message to gRPC CreateMessageResponse
func (m *MessageDTOMapper) ToCreateMessageResponse(msg domain.Message) *chatv1.CreateMessageResponse {
	return &chatv1.CreateMessageResponse{
		Message: &chatv1.Message{
			Id:         msg.ID,
			SenderId:   msg.SenderID,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  timestamppb.New(msg.CreatedAt),
		},
	}
}

// ToMessage converts domain Message to gRPC Message
func (m *MessageDTOMapper) ToMessage(msg domain.Message) *chatv1.Message {
	return &chatv1.Message{
		Id:         msg.ID,
		SenderId:   msg.SenderID,
		ReceiverId: msg.ReceiverID,
		Content:    msg.Content,
		CreatedAt:  timestamppb.New(msg.CreatedAt),
	}
}

// ToMessageList converts a slice of domain Messages to gRPC Messages
func (m *MessageDTOMapper) ToMessageList(messages []domain.Message) []*chatv1.Message {
	result := make([]*chatv1.Message, len(messages))
	for i, msg := range messages {
		result[i] = m.ToMessage(msg)
	}
	return result
}
