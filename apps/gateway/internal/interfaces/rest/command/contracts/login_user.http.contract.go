package contracts

import "github.com/gin-gonic/gin"

type LoginUserHTTPHandler interface {
	Mount(router *gin.RouterGroup)
}
