package query

import (
	"context"

	"golang-social-media/apps/gateway/internal/application/query/contracts"
	"golang-social-media/apps/gateway/internal/application/query/dto"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/logger"

	"github.com/rs/zerolog"
)

type getUserProfileQuery struct {
	authClient authProfileClient
	log        *zerolog.Logger
}

type authProfileClient interface {
	GetProfile(ctx context.Context, userID string) (auth.ProfileResponse, error)
}

func NewGetUserProfileQuery(client authProfileClient) contracts.GetUserProfileQuery {
	return &getUserProfileQuery{
		authClient: client,
		log:        logger.Component("gateway.query.get_user_profile"),
	}
}

func (q *getUserProfileQuery) Handle(ctx context.Context, userID string) (dto.UserProfileQueryResponse, error) {
	resp, err := q.authClient.GetProfile(ctx, userID)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to fetch profile from auth-service")
		return dto.UserProfileQueryResponse{}, err
	}

	q.log.Debug().
		Str("user_id", userID).
		Msg("fetched user profile")

	return dto.UserProfileQueryResponse{
		ID:    resp.ID,
		Email: resp.Email,
		Name:  resp.Name,
	}, nil
}
