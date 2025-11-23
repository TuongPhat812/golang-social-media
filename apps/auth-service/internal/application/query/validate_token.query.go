package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/redis"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
)

var _ contracts.ValidateTokenQuery = (*validateTokenQuery)(nil)

type validateTokenQuery struct {
	jwtService         *jwt.Service
	tokenBlacklistRepo *redis.TokenBlacklistRepository
	log                *zerolog.Logger
}

func NewValidateTokenQuery(
	jwtService *jwt.Service,
	tokenBlacklistRepo *redis.TokenBlacklistRepository,
) contracts.ValidateTokenQuery {
	return &validateTokenQuery{
		jwtService:         jwtService,
		tokenBlacklistRepo: tokenBlacklistRepo,
		log:                logger.Component("auth.query.validate_token"),
	}
}

func (q *validateTokenQuery) Execute(ctx context.Context, token string) (contracts.ValidateTokenQueryResponse, error) {
	// Validate token
	userID, err := q.jwtService.ValidateToken(token)
	if err != nil {
		q.log.Warn().
			Err(err).
			Msg("invalid token")
		return contracts.ValidateTokenQueryResponse{
			Valid:  false,
			UserID: "",
		}, nil // Return valid=false, not error
	}

	// Check if token is blacklisted
	tokenID := user.NewTokenID(token)
	isBlacklisted, err := q.tokenBlacklistRepo.IsBlacklisted(ctx, tokenID.String())
	if err != nil {
		q.log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to check token blacklist")
		return contracts.ValidateTokenQueryResponse{
			Valid:  false,
			UserID: "",
		}, nil // Return valid=false on error
	}

	if isBlacklisted {
		q.log.Warn().
			Str("user_id", userID).
			Msg("token is blacklisted")
		return contracts.ValidateTokenQueryResponse{
			Valid:  false,
			UserID: "",
		}, nil
	}

	return contracts.ValidateTokenQueryResponse{
		Valid:  true,
		UserID: userID,
	}, nil
}
