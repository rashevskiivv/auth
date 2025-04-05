package utils

import (
	"net/http"

	"github.com/rashevskiivv/auth/internal/entity"
	handlerAuth "github.com/rashevskiivv/auth/internal/handler/auth"

	"github.com/gin-gonic/gin"
)

// NotFound Not found page handler.
func NotFound(c *gin.Context) {
	b := entity.Response{
		Message: "Page not found",
		Errors:  "page not found",
	}
	c.JSON(http.StatusNotFound, b)
}

// HealthCheck Healthcheck page handler.
func HealthCheck(c *gin.Context) {
	b := entity.Response{
		Message: "SERVING",
	}
	c.JSON(http.StatusOK, b)
}

func TokenAuthMiddleware(handler handlerAuth.HandlerI) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.CheckTokenHandle(c)

		c.Next()
	}
}
