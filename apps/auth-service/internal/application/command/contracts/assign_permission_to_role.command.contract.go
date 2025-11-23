package contracts

import "context"

type AssignPermissionToRoleCommandRequest struct {
	RoleID       string
	PermissionID string
}

type AssignPermissionToRoleCommand interface {
	Execute(ctx context.Context, req AssignPermissionToRoleCommandRequest) error
}

