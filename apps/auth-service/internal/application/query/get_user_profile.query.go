package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/logger"
)

var _ contracts.GetUserProfileQuery = (*getUserProfileQuery)(nil)

type getUserProfileQuery struct {
	repo *memory.UserRepository
	log  *zerolog.Logger
}

func NewGetUserProfileHandler(repo *memory.UserRepository) contracts.GetUserProfileQuery {
	return &getUserProfileQuery{
		repo: repo,
		log:  logger.Component("auth.query.get_user_profile"),
	}
}

func (q *getUserProfileQuery) Execute(ctx context.Context, userID string) (auth.ProfileResponse, error) {
	user, err := q.repo.GetByID(userID)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to get user profile")
		// Transform error - transformer will handle not found cases
		return auth.ProfileResponse{}, err
	}

	q.log.Info().
		Str("user_id", userID).
		Msg("user profile retrieved")

	return auth.ProfileResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
