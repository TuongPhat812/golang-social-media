package contracts

import (
	"context"

	domain "golang-social-media/apps/gateway/internal/domain/message"
)

type CreateMessageCommand interface {
	Handle(ctx context.Context, senderID, receiverID, content string) (domain.Message, error)
}
