package query

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
	"golang-social-media/pkg/logger"
)

var _ contracts.CheckPermissionQuery = (*checkPermissionQuery)(nil)

type checkPermissionQuery struct {
	userRoleRepo        *postgres.UserRoleRepository
	rolePermissionRepo  *postgres.RolePermissionRepository
	permissionRepo      *postgres.PermissionRepository
	log                 *zerolog.Logger
}

func NewCheckPermissionQuery(
	userRoleRepo *postgres.UserRoleRepository,
	rolePermissionRepo *postgres.RolePermissionRepository,
	permissionRepo *postgres.PermissionRepository,
) contracts.CheckPermissionQuery {
	return &checkPermissionQuery{
		userRoleRepo:       userRoleRepo,
		rolePermissionRepo: rolePermissionRepo,
		permissionRepo:     permissionRepo,
		log:                logger.Component("auth.query.check_permission"),
	}
}

func (q *checkPermissionQuery) Execute(ctx context.Context, req contracts.CheckPermissionQueryRequest) (contracts.CheckPermissionQueryResponse, error) {
	// Get permission by resource and action
	perm, err := q.permissionRepo.GetByResourceAction(req.Resource, req.Action)
	if err != nil {
		q.log.Debug().
			Str("resource", req.Resource).
			Str("action", req.Action).
			Msg("permission not found")
		return contracts.CheckPermissionQueryResponse{
			HasPermission: false,
		}, nil
	}

	// Get user roles
	roleIDs, err := q.userRoleRepo.GetUserRoles(req.UserID)
	if err != nil {
		q.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Msg("failed to get user roles")
		return contracts.CheckPermissionQueryResponse{
			HasPermission: false,
		}, err
	}

	// Check if any of user's roles has this permission
	for _, roleID := range roleIDs {
		hasPermission, err := q.rolePermissionRepo.HasPermission(roleID, perm.ID)
		if err != nil {
			q.log.Warn().
				Err(err).
				Str("role_id", roleID).
				Str("permission_id", perm.ID).
				Msg("failed to check role permission")
			continue
		}
		if hasPermission {
			q.log.Debug().
				Str("user_id", req.UserID).
				Str("permission_id", perm.ID).
				Str("role_id", roleID).
				Msg("user has permission via role")
			return contracts.CheckPermissionQueryResponse{
				HasPermission: true,
			}, nil
		}
	}

	q.log.Debug().
		Str("user_id", req.UserID).
		Str("permission_id", perm.ID).
		Msg("user does not have permission")
	return contracts.CheckPermissionQueryResponse{
		HasPermission: false,
	}, nil
}

