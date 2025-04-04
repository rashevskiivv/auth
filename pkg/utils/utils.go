package utils

import (
	"net/http"
	"tax-auth/internal/entity"
	handlerAuth "tax-auth/internal/handler/auth"

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
		token := c.Request.Header.Get("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, entity.Response{Message: "Provided token is empty", Errors: "Provided token is empty"})
			return
		}

		requiredToken, err := handler.CheckTokenHandle(c, c.Request.Header.Get("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, entity.Response{Errors: err.Error()})
			return
		}

		if token != requiredToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.Response{Message: "Token is invalid", Errors: "Token is invalid"})
			return
		}

		c.Next()
	}
}
