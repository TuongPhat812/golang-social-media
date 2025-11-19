package command

import (
	"context"
	"strings"

	"golang-social-media/apps/auth-service/internal/application/command/contracts"
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
	repo            *memory.UserRepository
	userFactory     *factories.UserFactory
	eventDispatcher *event_dispatcher.Dispatcher
	log             *zerolog.Logger
}

func NewRegisterUserCommand(
	repo *memory.UserRepository,
	eventDispatcher *event_dispatcher.Dispatcher,
	idFn func() string,
) contracts.RegisterUserCommand {
	var factory *factories.UserFactory
	if idFn != nil {
		factory = factories.NewUserFactoryWithIDGenerator(idFn)
	} else {
		factory = factories.NewUserFactory()
	}

	return &registerUserCommand{
		repo:            repo,
		userFactory:     factory,
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

	// Persist to database
	if err := c.repo.Create(*userModel); err != nil {
		c.log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("failed to persist user")
		// Check if it's a duplicate email error
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			return auth.RegisterResponse{}, errors.NewConflictError(errors.CodeEmailAlreadyExists)
		}
		return auth.RegisterResponse{}, err
	}

	// Dispatch domain events AFTER successful persistence
	domainEvents := userModel.Events()
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
