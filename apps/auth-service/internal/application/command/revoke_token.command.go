package command

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
)

var _ contracts.RevokeTokenCommand = (*revokeTokenCommand)(nil)

type revokeTokenCommand struct {
	jwtService         *jwt.Service
	tokenBlacklistRepo *redis.TokenBlacklistRepository
	log                *zerolog.Logger
}

func NewRevokeTokenCommand(
	jwtService *jwt.Service,
	tokenBlacklistRepo *redis.TokenBlacklistRepository,
) contracts.RevokeTokenCommand {
	return &revokeTokenCommand{
		jwtService:         jwtService,
		tokenBlacklistRepo: tokenBlacklistRepo,
		log:                logger.Component("auth.command.revoke_token"),
	}
}

func (c *revokeTokenCommand) Execute(ctx context.Context, req contracts.RevokeTokenCommandRequest) error {
	// Validate token to extract user ID and expiration
	userID, err := c.jwtService.ValidateToken(req.Token)
	if err != nil {
		// Try refresh token validation
		userID, err = c.jwtService.ValidateRefreshToken(req.Token)
		if err != nil {
			c.log.Warn().
				Err(err).
				Msg("invalid token provided for revocation")
			return errors.NewUnauthorizedError("invalid token")
		}
		// It's a refresh token, use longer expiration
		tokenID := user.NewTokenID(req.Token)
		expiration := 8 * 24 * time.Hour // Refresh token expires in 7 days, add buffer
		if err := c.tokenBlacklistRepo.AddToken(ctx, tokenID.String(), expiration); err != nil {
			c.log.Error().
				Err(err).
				Str("user_id", userID).
				Msg("failed to revoke refresh token")
			return err
		}
		c.log.Info().
			Str("user_id", userID).
			Msg("refresh token revoked successfully")
		return nil
	}

	// It's an access token
	tokenID := user.NewTokenID(req.Token)
	expiration := 25 * time.Hour // Access token expires in 1 hour, add buffer
	if err := c.tokenBlacklistRepo.AddToken(ctx, tokenID.String(), expiration); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to revoke access token")
		return err
	}

	c.log.Info().
		Str("user_id", userID).
		Msg("access token revoked successfully")

	return nil
}

