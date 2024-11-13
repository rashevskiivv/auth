package auth

import (
	"log"
	"net/http"
	"tax-auth/internal/entity"
	"tax-auth/internal/usecase/auth"

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
