package contracts

import (
	"context"

	"golang-social-media/pkg/events"
)

// HandleChatCreatedCommand handles ChatCreated events
type HandleChatCreatedCommand interface {
	Execute(ctx context.Context, event events.ChatCreated) error
}

