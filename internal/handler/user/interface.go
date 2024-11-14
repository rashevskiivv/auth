package user

import "github.com/gin-gonic/gin"

type HandlerI interface {
	UpsertUserHandle(ctx *gin.Context)
	ReadUsersHandle(ctx *gin.Context)
	DeleteUsersHandle(ctx *gin.Context)
}
