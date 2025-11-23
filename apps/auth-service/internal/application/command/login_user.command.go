package command

import (
	"context"

	"golang-social-media/apps/auth-service/internal/infrastructure/jwt"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/pkg/contracts/auth"
)

type LoginUserHandler struct {
	repo      *memory.UserRepository
	jwtService *jwt.Service
}

func NewLoginUserHandler(repo *memory.UserRepository, jwtService *jwt.Service) *LoginUserHandler {
	return &LoginUserHandler{
		repo:       repo,
		jwtService: jwtService,
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

	// Generate JWT token pair (access + refresh)
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID)
	if err != nil {
		return auth.LoginResponse{}, err
	}

	return auth.LoginResponse{
		UserID:       user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}
