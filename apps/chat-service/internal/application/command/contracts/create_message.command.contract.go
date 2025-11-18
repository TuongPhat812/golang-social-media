package contracts

import (
	"context"

	"golang-social-media/apps/chat-service/internal/domain/message"
)

// CreateMessageCommand creates a new chat message
type CreateMessageCommand interface {
	Execute(ctx context.Context, req CreateMessageCommandRequest) (message.Message, error)
}

// CreateMessageCommandRequest represents the request for creating a message
type CreateMessageCommandRequest struct {
	SenderID   string
	ReceiverID string
	Content    string
}

