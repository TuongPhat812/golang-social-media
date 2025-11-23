package command

import (
	"context"

	"github.com/rs/zerolog"
	"golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/apps/auth-service/internal/domain/user_role"
	"golang-social-media/apps/auth-service/internal/infrastructure/persistence/postgres"
	"golang-social-media/pkg/logger"
)

var _ contracts.AssignRoleCommand = (*assignRoleCommand)(nil)

type assignRoleCommand struct {
	userRoleRepo *postgres.UserRoleRepository
	roleRepo     *postgres.RoleRepository
	log          *zerolog.Logger
}

func NewAssignRoleCommand(
	userRoleRepo *postgres.UserRoleRepository,
	roleRepo *postgres.RoleRepository,
) contracts.AssignRoleCommand {
	return &assignRoleCommand{
		userRoleRepo: userRoleRepo,
		roleRepo:     roleRepo,
		log:          logger.Component("auth.command.assign_role"),
	}
}

func (c *assignRoleCommand) Execute(ctx context.Context, req contracts.AssignRoleCommandRequest) error {
	// Verify role exists
	_, err := c.roleRepo.GetByID(req.RoleID)
	if err != nil {
		c.log.Error().
			Err(err).
			Str("role_id", req.RoleID).
			Msg("role not found")
		return err
	}

	// Create user role
	userRole := user_role.UserRole{
		UserID:    req.UserID,
		RoleID:    req.RoleID,
	}
	userRole.Assign()

	// Persist
	if err := c.userRoleRepo.Create(userRole); err != nil {
		c.log.Error().
			Err(err).
			Str("user_id", req.UserID).
			Str("role_id", req.RoleID).
			Msg("failed to assign role")
		return err
	}

	c.log.Info().
		Str("user_id", req.UserID).
		Str("role_id", req.RoleID).
		Msg("role assigned successfully")

	return nil
}

