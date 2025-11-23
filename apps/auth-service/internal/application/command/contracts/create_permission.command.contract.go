package contracts

import "context"

type CreatePermissionCommandRequest struct {
	Name     string
	Resource string
	Action   string
}

type CreatePermissionCommandResponse struct {
	ID       string
	Name     string
	Resource string
	Action   string
}

type CreatePermissionCommand interface {
	Execute(ctx context.Context, req CreatePermissionCommandRequest) (CreatePermissionCommandResponse, error)
}

