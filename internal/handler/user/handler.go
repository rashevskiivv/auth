package user

import (
	"errors"
	"io"
	"log"
	"net/http"
	"tax-auth/internal/entity"
	usecaseUser "tax-auth/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc usecaseUser.UseCaseI
}

func NewUserHandler(uc usecaseUser.UseCaseI) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) UpsertUserHandle(ctx *gin.Context) {
	var (
		input    entity.UpdateUsersInput
		response entity.Response
		err      error
	)

	if err = ctx.BindJSON(&input); err != nil {
		log.Println(err)
		response.Errors = err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = h.uc.UpdateUsers(ctx, input)
	if err != nil {
		log.Println(err)
		response.Errors = err.Error()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	response.Message = "Created"
	ctx.JSON(http.StatusCreated, response)
	return
}

func (h *Handler) ReadUsersHandle(ctx *gin.Context) {
	var (
		filter   entity.UserFilter
		response entity.Response
		err      error
	)

	err = ctx.ShouldBindJSON(&filter)
	if errors.Is(err, io.EOF) {
		err = nil
	}
	if err != nil {
		log.Println(err)
		response.Errors = err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	users, err := h.uc.GetUsers(ctx, entity.GetUsersInput{Filter: filter})
	if err != nil {
		log.Println(err)
		response.Errors = err.Error()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	if users == nil {
		response.Errors = "No users found"
		ctx.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}
	response.Data = users.Response
	ctx.JSON(http.StatusOK, response)
	return
}

func (h *Handler) DeleteUsersHandle(ctx *gin.Context) {
	var (
		filter   entity.UserFilter
		response entity.Response
		err      error
	)

	if err = ctx.BindJSON(&filter); err != nil {
		log.Println(err)
		response.Errors = err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = h.uc.DeleteUsers(ctx, entity.DeleteUsersInput{Filter: filter})
	if err != nil {
		log.Println(err)
		response.Errors = err.Error()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	response.Message = "Deleted"
	ctx.JSON(http.StatusOK, response)
	return

}
