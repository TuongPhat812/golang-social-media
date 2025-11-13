package contracts

import "github.com/gin-gonic/gin"

type GetUserProfileHTTPHandler interface {
	Mount(router *gin.RouterGroup)
}
