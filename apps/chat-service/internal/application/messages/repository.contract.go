package messages

import (
	"context"

	domain "golang-social-media/apps/chat-service/internal/domain/message"
)

type Repository interface {
	Create(ctx context.Context, msg *domain.Message) error
}
