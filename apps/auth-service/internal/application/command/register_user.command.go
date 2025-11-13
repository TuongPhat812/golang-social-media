package command

import (
	"context"

	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/apps/auth-service/internal/pkg/random"
	"golang-social-media/pkg/contracts/auth"
)

type RegisterUserHandler struct {
	repo *memory.UserRepository
	idFn func() string
}

func NewRegisterUserHandler(repo *memory.UserRepository, idFn func() string) *RegisterUserHandler {
	if idFn == nil {
		idFn = func() string {
			return "user-" + random.String(8)
		}
	}
	return &RegisterUserHandler{repo: repo, idFn: idFn}
}

func (h *RegisterUserHandler) Handle(ctx context.Context, req auth.RegisterRequest) (auth.RegisterResponse, error) {
	u := user.User{
		ID:       h.idFn(),
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}
	if err := h.repo.Create(u); err != nil {
		return auth.RegisterResponse{}, err
	}
	return auth.RegisterResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}, nil
}
