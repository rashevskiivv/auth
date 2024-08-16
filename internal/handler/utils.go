package handler

import (
	"net/http"
	"tax-auth/internal/entity"

	"github.com/gin-gonic/gin"
)

// NotFound Not found page handler.
func NotFound(c *gin.Context) {
	b := entity.Response{
		Data:    nil,
		Message: "Page not found",
		Errors:  "page not found",
	}
	c.JSON(http.StatusNotFound, b)
}

// HealthCheck Healthcheck page handler.
func HealthCheck(c *gin.Context) {
	b := entity.Response{
		Data:    nil,
		Message: "SERVING",
		Errors:  "",
	}
	c.JSON(http.StatusOK, b)
}
