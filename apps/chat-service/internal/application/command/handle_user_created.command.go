package command

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang-social-media/apps/chat-service/internal/application/command/contracts"
	"golang-social-media/apps/chat-service/internal/infrastructure/persistence"
	"golang-social-media/pkg/events"
	"golang-social-media/pkg/logger"
)

var _ contracts.HandleUserCreatedCommand = (*HandleUserCreatedCommandHandler)(nil)

type HandleUserCreatedCommandHandler struct {
	userRepo *persistence.UserRepository
	log      *zerolog.Logger
}

func NewHandleUserCreatedCommand(userRepo *persistence.UserRepository) *HandleUserCreatedCommandHandler {
	return &HandleUserCreatedCommandHandler{
		userRepo: userRepo,
		log:      logger.Component("chat.command.handle_user_created"),
	}
}

func (c *HandleUserCreatedCommandHandler) Execute(ctx context.Context, event events.UserCreated) error {
	userModel := persistence.UserModel{
		ID:        event.ID,
		Email:     event.Email,
		Name:      event.Name,
		CreatedAt: event.CreatedAt,
		UpdatedAt: time.Now().UTC(),
	}

	if err := c.userRepo.Upsert(ctx, userModel); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", event.ID).
			Str("email", event.Email).
			Msg("failed to replicate user")
		return err
	}

	c.log.Info().
		Str("user_id", event.ID).
		Str("email", event.Email).
		Str("name", event.Name).
		Msg("user replicated successfully")

	return nil
}

