package contracts

import "context"

type RevokeRoleCommandRequest struct {
	UserID string
	RoleID string
}

type RevokeRoleCommand interface {
	Execute(ctx context.Context, req RevokeRoleCommandRequest) error
}

