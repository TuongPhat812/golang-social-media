package query

import (
	"net/http"

	app "golang-social-media/apps/gateway/internal/application/query/contracts"
	httpcontracts "golang-social-media/apps/gateway/internal/interfaces/rest/query/contracts"
	"golang-social-media/pkg/logger"

	"github.com/gin-gonic/gin"
)

type getUserProfileHTTPHandler struct {
	query app.GetUserProfileQuery
}

func NewGetUserProfileHTTPHandler(query app.GetUserProfileQuery) httpcontracts.GetUserProfileHTTPHandler {
	return &getUserProfileHTTPHandler{query: query}
}

func (h *getUserProfileHTTPHandler) Mount(router *gin.RouterGroup) {
	router.GET("/profile/:id", h.handle)
}

func (h *getUserProfileHTTPHandler) handle(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	resp, err := h.query.Handle(c.Request.Context(), userID)
	if err != nil {
		logger.Component("gateway.http.get_user_profile").
			Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to get user profile")
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
