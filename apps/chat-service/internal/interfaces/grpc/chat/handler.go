package chat

import (
	"context"

	bootstrap "golang-social-media/apps/chat-service/internal/infrastructure/bootstrap"
	commandcontracts "golang-social-media/apps/chat-service/internal/application/command/contracts"
	chatv1 "golang-social-media/pkg/gen/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	createMessageCmd commandcontracts.CreateMessageCommand
	chatv1.UnimplementedChatServiceServer
}

func NewHandler(deps *bootstrap.Dependencies) *Handler {
	return &Handler{
		createMessageCmd: deps.CreateMessageCmd,
	}
}

func (h *Handler) CreateMessage(ctx context.Context, req *chatv1.CreateMessageRequest) (*chatv1.CreateMessageResponse, error) {
	msg, err := h.createMessageCmd.Execute(ctx, commandcontracts.CreateMessageCommandRequest{
		SenderID:   req.GetSenderId(),
		ReceiverID: req.GetReceiverId(),
		Content:    req.GetContent(),
	})
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
