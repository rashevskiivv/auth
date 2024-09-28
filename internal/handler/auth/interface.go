package auth

import "github.com/gin-gonic/gin"

type HandlerI interface {
	RegisterUserHandle(ctx *gin.Context)
	AuthenticateUserHandle(ctx *gin.Context)
}
