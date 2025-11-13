package contracts

import (
	"context"

	"golang-social-media/apps/gateway/internal/application/command/dto"
	domain "golang-social-media/apps/gateway/internal/domain/user"
)

type RegisterUserCommand interface {
	Handle(ctx context.Context, req dto.RegisterUserCommandRequest) (domain.User, error)
}
