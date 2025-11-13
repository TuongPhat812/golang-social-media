package command

import (
	"context"

	"golang-social-media/apps/gateway/internal/application/command/contracts"
	"golang-social-media/apps/gateway/internal/application/command/dto"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

type loginUserCommand struct {
	authClient authLoginClient
	log        *zerolog.Logger
}

type authLoginClient interface {
	Login(ctx context.Context, req auth.LoginRequest) (auth.LoginResponse, error)
}

func NewLoginUserCommand(client authLoginClient) contracts.LoginUserCommand {
	return &loginUserCommand{
		authClient: client,
		log:        logger.Component("gateway.command.login_user"),
	}
}

func (c *loginUserCommand) Handle(ctx context.Context, req dto.LoginUserCommandRequest) (dto.LoginUserCommandResponse, error) {
	resp, err := c.authClient.Login(ctx, auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("failed to login user via auth-service")
		return dto.LoginUserCommandResponse{}, err
	}

	c.log.Info().
		Str("user_id", resp.UserID).
		Msg("user login successful")

	return dto.LoginUserCommandResponse{
		UserID: resp.UserID,
		Token:  resp.Token,
	}, nil
}
