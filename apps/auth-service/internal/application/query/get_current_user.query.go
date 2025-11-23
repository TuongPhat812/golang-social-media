package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/errors"
	"golang-social-media/pkg/logger"
)

var _ contracts.GetCurrentUserQuery = (*getCurrentUserQuery)(nil)

type getCurrentUserQuery struct {
	userRepo *memory.UserRepository
	log      *zerolog.Logger
}

func NewGetCurrentUserQuery(userRepo *memory.UserRepository) contracts.GetCurrentUserQuery {
	return &getCurrentUserQuery{
		userRepo: userRepo,
		log:      logger.Component("auth.query.get_current_user"),
	}
}

func (q *getCurrentUserQuery) Execute(ctx context.Context, userID string) (contracts.GetCurrentUserQueryResponse, error) {
	user, err := q.userRepo.GetByID(userID)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to get user")
		return contracts.GetCurrentUserQueryResponse{}, errors.NewNotFoundError(errors.CodeUserNotFound)
	}

	return contracts.GetCurrentUserQueryResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
