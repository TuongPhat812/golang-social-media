package chat

import (
	"context"

	"github.com/myself/golang-social-media/apps/chat-service/internal/application/messages"
	chatcontract "github.com/myself/golang-social-media/pkg/contracts/chat"
)

type Handler struct {
	service messages.Service
	chatcontract.UnimplementedChatServiceServer
}

func NewHandler(service messages.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateMessage(ctx context.Context, req *chatcontract.CreateMessageRequest) (*chatcontract.CreateMessageResponse, error) {
	msg, err := h.service.CreateMessage(ctx, req.SenderID, req.ReceiverID, req.Content)
	if err != nil {
		return nil, err
	}

	return &chatcontract.CreateMessageResponse{
		Message: chatcontract.Message{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  msg.CreatedAt,
		},
	}, nil
}
