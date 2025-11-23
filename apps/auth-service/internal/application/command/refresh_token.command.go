package command

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
)

var _ contracts.RefreshTokenCommand = (*refreshTokenCommand)(nil)

type refreshTokenCommand struct {
	userRepo           *memory.UserRepository
	jwtService         *jwt.Service
	tokenBlacklistRepo *redis.TokenBlacklistRepository
	log                *zerolog.Logger
}

func NewRefreshTokenCommand(
	userRepo *memory.UserRepository,
	jwtService *jwt.Service,
	tokenBlacklistRepo *redis.TokenBlacklistRepository,
) contracts.RefreshTokenCommand {
	return &refreshTokenCommand{
		userRepo:           userRepo,
		jwtService:         jwtService,
		tokenBlacklistRepo: tokenBlacklistRepo,
		log:                logger.Component("auth.command.refresh_token"),
	}
}

func (c *refreshTokenCommand) Execute(ctx context.Context, req contracts.RefreshTokenCommandRequest) (contracts.RefreshTokenCommandResponse, error) {
	// Validate refresh token
	userID, err := c.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.log.Warn().
			Err(err).
			Msg("invalid refresh token")
		return contracts.RefreshTokenCommandResponse{}, errors.NewUnauthorizedError("invalid refresh token")
	}

	// Check if refresh token is blacklisted
	tokenID := user.NewTokenID(req.RefreshToken)
	isBlacklisted, err := c.tokenBlacklistRepo.IsBlacklisted(ctx, tokenID.String())
	if err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to check token blacklist")
		return contracts.RefreshTokenCommandResponse{}, err
	}
	if isBlacklisted {
		c.log.Warn().
			Str("user_id", userID).
			Msg("refresh token is blacklisted")
		return contracts.RefreshTokenCommandResponse{}, errors.NewUnauthorizedError("token is blacklisted")
	}

	// Verify user exists
	_, err = c.userRepo.GetByID(userID)
	if err != nil {
		c.log.Warn().
			Err(err).
			Str("user_id", userID).
			Msg("user not found")
		return contracts.RefreshTokenCommandResponse{}, errors.NewUnauthorizedError("user not found")
	}

	// Generate new token pair
	tokenPair, err := c.jwtService.GenerateTokenPair(userID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to generate token pair")
		return contracts.RefreshTokenCommandResponse{}, err
	}

	// Blacklist old refresh token
	oldTokenID := user.NewTokenID(req.RefreshToken)
	// Refresh token expires in 7 days, add some buffer
	if err := c.tokenBlacklistRepo.AddToken(ctx, oldTokenID.String(), 8*24*time.Hour); err != nil {
		c.log.Warn().
			Err(err).
			Str("user_id", userID).
			Msg("failed to blacklist old refresh token")
		// Don't fail the request, just log warning
	}

	c.log.Info().
		Str("user_id", userID).
		Msg("token refreshed successfully")

	return contracts.RefreshTokenCommandResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}
