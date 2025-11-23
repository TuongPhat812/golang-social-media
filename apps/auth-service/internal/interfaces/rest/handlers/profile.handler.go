package handlers

import (
	"net/http"

	commandcontracts "golang-social-media/apps/auth-service/internal/application/command/contracts"
	querycontracts "golang-social-media/apps/auth-service/internal/application/query/contracts"
	"golang-social-media/pkg/contracts/auth"
	"golang-social-media/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ProfileHandler handles profile-related endpoints
type ProfileHandler struct {
	updateProfile  commandcontracts.UpdateProfileCommand
	getProfile     querycontracts.GetUserProfileQuery
	getCurrentUser querycontracts.GetCurrentUserQuery
}

// NewProfileHandler creates a new ProfileHandler
func NewProfileHandler(
	updateProfile commandcontracts.UpdateProfileCommand,
	getProfile querycontracts.GetUserProfileQuery,
	getCurrentUser querycontracts.GetCurrentUserQuery,
) *ProfileHandler {
	return &ProfileHandler{
		updateProfile:  updateProfile,
		getProfile:     getProfile,
		getCurrentUser: getCurrentUser,
	}
}

// Mount mounts profile routes to the router group
func (h *ProfileHandler) Mount(group *gin.RouterGroup) {
	// Public route
	group.GET("/profile/:id", h.getProfileByID)
}

// MountProtected mounts protected profile routes (require JWT middleware)
func (h *ProfileHandler) MountProtected(group *gin.RouterGroup) {
	group.GET("/me", h.getCurrentUserProfile)
	group.PUT("/profile", h.updateUserProfile)
}

// getProfileByID handles GET /auth/profile/:id
func (h *ProfileHandler) getProfileByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.getProfile.Execute(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// getCurrentUserProfile handles GET /auth/me
func (h *ProfileHandler) getCurrentUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("User ID not found in context"))
		return
	}

	resp, err := h.getCurrentUser.Execute(c.Request.Context(), userID.(string))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, auth.ProfileResponse{
		ID:    resp.ID,
		Email: resp.Email,
		Name:  resp.Name,
	})
}

// updateUserProfile handles PUT /auth/profile
func (h *ProfileHandler) updateUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.NewUnauthorizedError("User ID not found in context"))
		return
	}

	var req auth.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewInvalidRequestError("Invalid request body"))
		return
	}

	resp, err := h.updateProfile.Execute(c.Request.Context(), commandcontracts.UpdateProfileCommandRequest{
		UserID: userID.(string),
		Name:   req.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, auth.UpdateProfileResponse{
		ID:    resp.ID,
		Email: resp.Email,
		Name:  resp.Name,
	})
}

