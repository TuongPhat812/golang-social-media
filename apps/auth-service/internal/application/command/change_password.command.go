package command

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
)

var _ contracts.ChangePasswordCommand = (*changePasswordCommand)(nil)

type changePasswordCommand struct {
	userRepo        *memory.UserRepository
	eventDispatcher eventDispatcher
	log             *zerolog.Logger
}

type eventDispatcher interface {
	Dispatch(ctx context.Context, event interface{}) error
}

func NewChangePasswordCommand(
	userRepo *memory.UserRepository,
	eventDispatcher eventDispatcher,
) contracts.ChangePasswordCommand {
	return &changePasswordCommand{
		userRepo:        userRepo,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("auth.command.change_password"),
	}
}

func (c *changePasswordCommand) Execute(ctx context.Context, req contracts.ChangePasswordCommandRequest) error {
	// Get user from repository
	userEntity, err := c.userRepo.GetByID(req.UserID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to get user")
		return err
	}

	// Verify current password
	if userEntity.Password != req.CurrentPassword {
		c.log.Warn().
			Str("user_id", req.UserID).
			Msg("invalid current password")
		return errors.NewValidationError(errors.CodeInvalidCredentials, nil)
	}

	// Validate new password using domain method
	if err := userEntity.ValidatePassword(req.NewPassword); err != nil {
		c.log.Warn().
			Err(err).
			Str("user_id", req.UserID).
			Msg("invalid new password")
		return err
	}

	// Change password using domain method
	userEntity.ChangePassword(req.NewPassword)

	// Persist changes
	if err := c.userRepo.Update(userEntity); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to update user password")
		return err
	}

	// Dispatch domain events
	domainEvents := userEntity.Events()
	userEntity.ClearEvents()

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("user_id", req.UserID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("user_id", req.UserID).
		Msg("password changed successfully")

	return nil
}
