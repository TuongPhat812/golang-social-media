package contracts

import "context"

type GetUserRolesQueryResponse struct {
	UserID string
	RoleIDs []string
	RoleNames []string
}

type GetUserRolesQuery interface {
	Execute(ctx context.Context, userID string) (GetUserRolesQueryResponse, error)
}

