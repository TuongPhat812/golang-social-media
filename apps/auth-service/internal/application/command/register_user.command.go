package command

import (
	"context"
	"time"

	"golang-social-media/apps/auth-service/internal/domain/user"
	"golang-social-media/apps/auth-service/internal/infrastructure/eventbus"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/memory"
	"golang-social-media/apps/auth-service/internal/pkg/random"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/events"
)

type RegisterUserHandler struct {
	repo      *memory.UserRepository
	publisher *eventbus.KafkaPublisher
	idFn      func() string
}

func NewRegisterUserHandler(repo *memory.UserRepository, publisher *eventbus.KafkaPublisher, idFn func() string) *RegisterUserHandler {
	if idFn == nil {
		idFn = func() string {
			return "user-" + random.String(8)
		}
	}
	return &RegisterUserHandler{repo: repo, publisher: publisher, idFn: idFn}
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

	if h.publisher != nil {
		event := events.UserCreated{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			CreatedAt: time.Now().UTC(),
		}
		if err := h.publisher.PublishUserCreated(ctx, event); err != nil {
			return auth.RegisterResponse{}, err
		}
	}

	return auth.RegisterResponse{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}, nil
}
