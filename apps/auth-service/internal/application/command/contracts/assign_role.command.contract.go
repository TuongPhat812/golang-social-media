package contracts

import "context"

type AssignRoleCommandRequest struct {
	UserID string
	RoleID string
}

type AssignRoleCommand interface {
	Execute(ctx context.Context, req AssignRoleCommandRequest) error
}

