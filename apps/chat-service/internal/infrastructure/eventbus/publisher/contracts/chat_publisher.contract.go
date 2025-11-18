package contracts

import (
	"context"

	"golang-social-media/pkg/events"
)

// ChatPublisher publishes chat-related events
type ChatPublisher interface {
	PublishChatCreated(ctx context.Context, event events.ChatCreated) error
	Close() error
}
