package contracts

import "context"

type CreateRoleCommandRequest struct {
	Name        string
	Description string
}

type CreateRoleCommandResponse struct {
	ID          string
	Name        string
	Description string
}

type CreateRoleCommand interface {
	Execute(ctx context.Context, req CreateRoleCommandRequest) (CreateRoleCommandResponse, error)
}

