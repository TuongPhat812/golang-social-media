package query

import (
	"context"

	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/contracts/auth"
)

type GetUserProfileHandler struct {
	repo *memory.UserRepository
}

func NewGetUserProfileHandler(repo *memory.UserRepository) *GetUserProfileHandler {
	return &GetUserProfileHandler{repo: repo}
}

func (h *GetUserProfileHandler) Handle(ctx context.Context, userID string) (auth.ProfileResponse, error) {
	user, err := h.repo.GetByID(userID)
	if err != nil {
		return auth.ProfileResponse{}, err
	}
	return auth.ProfileResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
