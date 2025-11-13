package contracts

import "github.com/gin-gonic/gin"

type CreateMessageHTTPHandler interface {
	Mount(router *gin.RouterGroup)
}
