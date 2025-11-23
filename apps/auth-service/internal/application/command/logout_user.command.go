package command

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"
	"golang-social-media/pkg/logger"
)

var _ contracts.LogoutUserCommand = (*logoutUserCommand)(nil)

type logoutUserCommand struct {
	tokenBlacklistRepo *redis.TokenBlacklistRepository
	log                *zerolog.Logger
}

func NewLogoutUserCommand(tokenBlacklistRepo *redis.TokenBlacklistRepository) contracts.LogoutUserCommand {
	return &logoutUserCommand{
		tokenBlacklistRepo: tokenBlacklistRepo,
		log:                logger.Component("auth.command.logout_user"),
	}
}

func (c *logoutUserCommand) Execute(ctx context.Context, req contracts.LogoutUserCommandRequest) error {
	// Create token ID from token string
	tokenID := user.NewTokenID(req.Token)

	// Calculate expiration time (assume token expires in 24h, add some buffer)
	expiration := 25 * time.Hour

	// Add token to blacklist
	if err := c.tokenBlacklistRepo.AddToken(ctx, tokenID.String(), expiration); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to blacklist token")
		return err
	}

	c.log.Info().
		Str("user_id", req.UserID).
		Msg("user logged out successfully")

	return nil
}
