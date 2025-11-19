package command

import (
	"context"
	"strings"

	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	event_dispatcher "golang-social-media/apps/auth-service/internal/application/event_dispatcher"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/apps/auth-service/internal/pkg/random"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

var _ contracts.RegisterUserCommand = (*registerUserCommand)(nil)

type registerUserCommand struct {
	repo            *memory.UserRepository
	eventDispatcher *event_dispatcher.Dispatcher
	idFn            func() string
	log             *zerolog.Logger
}

func NewRegisterUserCommand(
	repo *memory.UserRepository,
	eventDispatcher *event_dispatcher.Dispatcher,
	idFn func() string,
) contracts.RegisterUserCommand {
	if idFn == nil {
		idFn = func() string {
			return "user-" + random.String(8)
		}
	}
	return &registerUserCommand{
		repo:            repo,
		eventDispatcher: eventDispatcher,
		idFn:            idFn,
		log:             logger.Component("auth.command.register_user"),
	}
}

func (c *registerUserCommand) Execute(ctx context.Context, req auth.RegisterRequest) (auth.RegisterResponse, error) {
	userModel := user.User{
		ID:       c.idFn(),
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	// Validate business rules before persisting or publishing
	if err := userModel.Validate(); err != nil {
		c.log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("user validation failed")
		return auth.RegisterResponse{}, err
	}

	// Domain logic: create user (this adds domain events internally)
	userModel.Create()

	// Persist to database
	if err := c.repo.Create(userModel); err != nil {
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
