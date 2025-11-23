package command

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
)

var _ contracts.UpdateProfileCommand = (*updateProfileCommand)(nil)

type updateProfileCommand struct {
	userRepo        *memory.UserRepository
	userFactory     factories.UserFactory
	eventDispatcher eventDispatcher
	log             *zerolog.Logger
}

type eventDispatcher interface {
	Dispatch(ctx context.Context, event interface{}) error
}

func NewUpdateProfileCommand(
	userRepo *memory.UserRepository,
	userFactory factories.UserFactory,
	eventDispatcher eventDispatcher,
) contracts.UpdateProfileCommand {
	return &updateProfileCommand{
		userRepo:        userRepo,
		userFactory:     userFactory,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("auth.command.update_profile"),
	}
}

func (c *updateProfileCommand) Execute(ctx context.Context, req contracts.UpdateProfileCommandRequest) (contracts.UpdateProfileCommandResponse, error) {
	// Get user from repository
	user, err := c.userRepo.GetByID(req.UserID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to get user")
		return contracts.UpdateProfileCommandResponse{}, err
	}

	// Validate new name
	if req.Name != "" {
		// Domain validation
		if len(req.Name) < 1 {
			return contracts.UpdateProfileCommandResponse{}, errors.NewValidationError(errors.CodeNameRequired, nil)
		}

		// Update profile using domain method
		user.UpdateProfile(req.Name)
	} else {
		// No changes
		return contracts.UpdateProfileCommandResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		}, nil
	}

	// Persist changes
	if err := c.userRepo.Update(user); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", user.ID).
			Msg("failed to update user")
		return contracts.UpdateProfileCommandResponse{}, err
	}

	// Dispatch domain events
	domainEvents := user.Events()
	user.ClearEvents()

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("user_id", user.ID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("user_id", user.ID).
		Str("new_name", req.Name).
		Msg("user profile updated")

	return contracts.UpdateProfileCommandResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
