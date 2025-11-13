package contracts

import (
	"context"

	"golang-social-media/apps/gateway/internal/application/command/dto"
)

type LoginUserCommand interface {
	Handle(ctx context.Context, req dto.LoginUserCommandRequest) (dto.LoginUserCommandResponse, error)
}
