package command

import (
	"context"
	"strings"

	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/application/repository"
	"golang-social-media/apps/auth-service/internal/application/unit_of_work"
	event_dispatcher "golang-social-media/apps/auth-service/internal/application/event_dispatcher"
	"golang-social-media/apps/auth-service/internal/domain/factories"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

var _ contracts.RegisterUserCommand = (*registerUserCommand)(nil)

type registerUserCommand struct {
	userRepo        repository.UserRepository
	uowFactory      unit_of_work.Factory
	userFactory     factories.UserFactory
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewRegisterUserCommand(
	userRepo repository.UserRepository,
	userFactory factories.UserFactory,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.RegisterUserCommand {
	return &registerUserCommand{
		userRepo:        userRepo,
		userFactory:     userFactory,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("auth.command.register_user"),
	}
}

// NewRegisterUserCommandWithUoW creates a new command with Unit of Work support
func NewRegisterUserCommandWithUoW(
	uowFactory unit_of_work.Factory,
	userFactory factories.UserFactory,
	eventDispatcher *event_dispatcher.Dispatcher,
) contracts.RegisterUserCommand {
	return &registerUserCommand{
		uowFactory:      uowFactory,
		userFactory:     userFactory,
		eventDispatcher: eventDispatcher,
		log:             logger.Component("auth.command.register_user"),
	}
}

func (c *registerUserCommand) Execute(ctx context.Context, req auth.RegisterRequest) (auth.RegisterResponse, error) {
	// Use factory to create user
	userModel, err := c.userFactory.CreateUser(req.Email, req.Password, req.Name)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("failed to create user using factory")
		return auth.RegisterResponse{}, err
	}

	// Get domain events BEFORE persisting
	domainEvents := userModel.Events()

	// If we have Unit of Work, use transactional approach
	if c.uowFactory != nil {
		// Create unit of work
		uow, err := c.uowFactory.New(ctx)
		if err != nil {
			c.log.Error().
				Err(err).
				Str("email", req.Email).
				Msg("failed to create unit of work")
			return auth.RegisterResponse{}, err
		}
		defer uow.Rollback() // Ensure rollback if commit fails

		// Persist user within transaction
		if err := uow.Users().Create(*userModel); err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
				return auth.RegisterResponse{}, errors.NewConflictError(errors.CodeEmailAlreadyExists)
			}
			c.log.Error().
				Err(err).
				Str("email", req.Email).
				Msg("failed to persist user")
			return auth.RegisterResponse{}, err
		}

		// Save events to outbox and event store within the same transaction
		events := make([]interface{}, len(domainEvents))
		for i, event := range domainEvents {
			events[i] = event
		}
		if err := uow.SaveEvents(ctx, events); err != nil {
			c.log.Error().
				Err(err).
				Str("email", req.Email).
				Msg("failed to save events")
			return auth.RegisterResponse{}, err
		}

		// Commit transaction (includes outbox and event store writes)
		if err := uow.Commit(); err != nil {
			c.log.Error().
				Err(err).
				Str("email", req.Email).
				Msg("failed to commit transaction")
			return auth.RegisterResponse{}, err
		}

		// Clear events after successful persistence
		userModel.ClearEvents()

		c.log.Info().
			Str("user_id", userModel.ID).
			Str("email", userModel.Email).
			Int("event_count", len(domainEvents)).
			Msg("user created with events saved to outbox and event store")

		return auth.RegisterResponse{
			ID:    userModel.ID,
			Email: userModel.Email,
			Name:  userModel.Name,
		}, nil
	}

	// Fallback: Use repository directly (for memory repo or when UoW not available)
	if err := c.userRepo.Create(*userModel); err != nil {
		c.log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("failed to persist user")
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			return auth.RegisterResponse{}, errors.NewConflictError(errors.CodeEmailAlreadyExists)
		}
		return auth.RegisterResponse{}, err
	}

	// Dispatch domain events AFTER successful persistence
	userModel.ClearEvents() // Clear events after dispatch

	for _, domainEvent := range domainEvents {
		if err := c.eventDispatcher.Dispatch(ctx, domainEvent); err != nil {
			// Log error but don't fail the command
			// Events can be retried via outbox pattern in production
			c.log.Error().
				Err(err).
				Str("event_type", domainEvent.Type()).
				Str("user_id", userModel.ID).
				Msg("failed to dispatch domain event")
		}
	}

	c.log.Info().
		Str("user_id", userModel.ID).
		Str("email", userModel.Email).
		Msg("user created")

	return auth.RegisterResponse{
		ID:    userModel.ID,
		Email: userModel.Email,
		Name:  userModel.Name,
	}, nil
}
