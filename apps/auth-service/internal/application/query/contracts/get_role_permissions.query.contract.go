package contracts

import "context"

type GetRolePermissionsQueryResponse struct {
	RoleID         string
	PermissionIDs  []string
	Permissions    []PermissionInfo
}

type PermissionInfo struct {
	ID       string
	Name     string
	Resource string
	Action   string
}

type GetRolePermissionsQuery interface {
	Execute(ctx context.Context, roleID string) (GetRolePermissionsQueryResponse, error)
}

