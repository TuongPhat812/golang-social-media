package command

import (
	"context"

	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/apps/auth-service/internal/pkg/random"
	"golang-social-media/pkg/contracts/auth"
)

type LoginUserHandler struct {
	repo       *memory.UserRepository
	tokenStore *TokenStore
}

func NewLoginUserHandler(repo *memory.UserRepository, tokenStore *TokenStore) *LoginUserHandler {
	return &LoginUserHandler{
		repo:       repo,
		tokenStore: tokenStore,
	}
}

func (h *LoginUserHandler) Handle(ctx context.Context, req auth.LoginRequest) (auth.LoginResponse, error) {
	user, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		return auth.LoginResponse{}, err
	}
	if user.Password != req.Password {
		return auth.LoginResponse{}, memory.ErrInvalidAuth
	}
	token := random.String(16)
	h.tokenStore.Save(token, user.ID)

	return auth.LoginResponse{
		UserID: user.ID,
		Token:  token,
	}, nil
}
