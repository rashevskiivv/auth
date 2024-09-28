package user

import "github.com/gin-gonic/gin"

type HandlerI interface {
	InsertUserHandle(ctx *gin.Context)
	ReadUsersHandle(ctx *gin.Context)
}
