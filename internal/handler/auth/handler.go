package auth

import (
	"log"
	"net/http"

	"github.com/rashevskiivv/auth/internal/entity"
	"github.com/rashevskiivv/auth/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc auth.UseCaseI
}

func NewAuthHandler(uc auth.UseCaseI) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) RegisterUserHandle(ctx *gin.Context) {
	var (
		err error
	)
	var request entity.User
	err = ctx.BindJSON(&request)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Errors: err.Error(),
		})
		return
	}

	registerInput := entity.RegisterInput{User: request}
	token, err := h.uc.RegisterUser(ctx, registerInput)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, entity.Response{
			Errors: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, entity.Response{
		Data:    token,
		Message: "User registered. Access token in data.",
	})
	return
}

func (h *Handler) AuthenticateUserHandle(ctx *gin.Context) {
	var request entity.User
	err := ctx.BindJSON(&request)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, entity.Response{
			Errors: err.Error(),
		})
		return
	}

	loginInput := entity.AuthenticateInput{User: request}
	token, err := h.uc.AuthenticateUser(ctx, loginInput)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, entity.Response{Errors: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, entity.Response{
		Data:    token,
		Message: "User logged in. Access token in data.",
	})
	return
}

func (h *Handler) CheckTokenHandle(ctx *gin.Context) {
	id := ctx.Request.Header.Get("id")
	token := ctx.Request.Header.Get("token")
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.Response{Message: "Provided id is empty", Errors: "Provided id is empty"})
		return
	}
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.Response{Message: "Provided token is empty", Errors: "Provided token is empty"})
		return
	}

	requiredToken, err := h.uc.CheckToken(ctx, entity.CheckTokenInput{UserID: id})
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, entity.Response{Errors: err.Error()})
		return
	}
	if requiredToken.Token.Token != token {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, entity.Response{Message: "Token is invalid", Errors: "Token is invalid"})
		return
	}
	return
}
