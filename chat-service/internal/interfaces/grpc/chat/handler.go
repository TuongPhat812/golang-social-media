package chat

import (
	"context"

	"github.com/myself/golang-social-media/chat-service/internal/application/messages"
	chatcontract "github.com/myself/golang-social-media/common/contracts/chat"
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

	return &chatcontract.CreateMessageResponse{Message: msg}, nil
}
