package contracts

import (
	"context"

	"golang-social-media/apps/gateway/internal/application/query/dto"
)

type GetUserProfileQuery interface {
	Handle(ctx context.Context, userID string) (dto.UserProfileQueryResponse, error)
}
