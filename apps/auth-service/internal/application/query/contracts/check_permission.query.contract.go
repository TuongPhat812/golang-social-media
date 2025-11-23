package contracts

import "context"

type CheckPermissionQueryRequest struct {
	UserID   string
	Resource string
	Action   string
}

type CheckPermissionQueryResponse struct {
	HasPermission bool
}

type CheckPermissionQuery interface {
	Execute(ctx context.Context, req CheckPermissionQueryRequest) (CheckPermissionQueryResponse, error)
}

