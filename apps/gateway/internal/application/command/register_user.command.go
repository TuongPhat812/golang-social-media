package command

import (
	"context"

	"golang-social-media/apps/gateway/internal/application/command/contracts"
	"golang-social-media/apps/gateway/internal/application/command/dto"
	domain "golang-social-media/apps/gateway/internal/domain/user"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

type registerUserCommand struct {
	authClient authCommandClient
	log        *zerolog.Logger
}

type authCommandClient interface {
	Register(ctx context.Context, req auth.RegisterRequest) (auth.RegisterResponse, error)
}

func NewRegisterUserCommand(client authCommandClient) contracts.RegisterUserCommand {
	return &registerUserCommand{
		authClient: client,
		log:        logger.Component("gateway.command.register_user"),
	}
}

func (c *registerUserCommand) Handle(ctx context.Context, req dto.RegisterUserCommandRequest) (domain.User, error) {
	resp, err := c.authClient.Register(ctx, auth.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		c.log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("failed to register user via auth-service")
		return domain.User{}, err
	}

	user := domain.User{
		ID:    resp.ID,
		Email: resp.Email,
		Name:  resp.Name,
	}

	c.log.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Msg("user registered successfully")

	return user, nil
}
