package handlers

import (
	"net/http"

	commandcontracts "golang-social-media/apps/auth-service/internal/application/command/contracts"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/errors"

	"github.com/gin-gonic/gin"
)

// PasswordHandler handles password-related endpoints
type PasswordHandler struct {
	changePassword commandcontracts.ChangePasswordCommand
}

// NewPasswordHandler creates a new PasswordHandler
func NewPasswordHandler(changePassword commandcontracts.ChangePasswordCommand) *PasswordHandler {
	return &PasswordHandler{
		changePassword: changePassword,
	}
}

// Mount mounts password routes to the router group
func (h *PasswordHandler) Mount(group *gin.RouterGroup) {
	group.POST("/change-password", h.changeUserPassword)
}

// changeUserPassword handles POST /auth/change-password
func (h *PasswordHandler) changeUserPassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("User ID not found in context"))
		return
	}

	var req auth.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidRequestError("Invalid request body"))
		return
	}

	err := h.changePassword.Execute(c.Request.Context(), commandcontracts.ChangePasswordCommandRequest{
		UserID:          userID.(string),
		CurrentPassword: req.OldPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

