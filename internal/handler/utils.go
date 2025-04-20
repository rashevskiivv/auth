package handler

import (
	"net/http"

	"github.com/rashevskiivv/auth/internal/entity"
	handlerAuth "github.com/rashevskiivv/auth/internal/handler/auth"

	"github.com/gin-gonic/gin"
)

// NotFound Not found page handler.
func NotFound(ctx *gin.Context) {
	b := entity.Response{
		Message: "Page not found",
		Errors:  "page not found",
	}
	ctx.JSON(http.StatusNotFound, b)
}

// HealthCheck Healthcheck page handler.
func HealthCheck(ctx *gin.Context) {
	b := entity.Response{
		Message: "SERVING",
	}
	ctx.JSON(http.StatusOK, b)
}

func TokenAuthMiddleware(handler handlerAuth.HandlerI) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler.CheckTokenHandle(ctx)

		ctx.Next()
	}
}
