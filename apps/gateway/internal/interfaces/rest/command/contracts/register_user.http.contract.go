package contracts

import "github.com/gin-gonic/gin"

type RegisterUserHTTPHandler interface {
	Mount(router *gin.RouterGroup)
}
