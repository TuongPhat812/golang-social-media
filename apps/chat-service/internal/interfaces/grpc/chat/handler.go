package chat

import (
	"context"

	"github.com/myself/golang-social-media/apps/chat-service/internal/application/messages"
	chatv1 "github.com/myself/golang-social-media/pkg/gen/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	service messages.Service
	chatv1.UnimplementedChatServiceServer
}

func NewHandler(service messages.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateMessage(ctx context.Context, req *chatv1.CreateMessageRequest) (*chatv1.CreateMessageResponse, error) {
	msg, err := h.service.CreateMessage(ctx, req.GetSenderId(), req.GetReceiverId(), req.GetContent())
	if err != nil {
		return nil, err
	}

	return &chatv1.CreateMessageResponse{
		Message: &chatv1.Message{
			Id:         msg.ID,
			SenderId:   msg.SenderID,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  timestamppb.New(msg.CreatedAt),
		},
	}, nil
}
