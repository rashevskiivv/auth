package user

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/rashevskiivv/auth/internal/entity"
	usecaseUser "github.com/rashevskiivv/auth/internal/usecase/user"

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
	log.Println("handle UpsertUserHandle started")
	defer log.Println("handle UpsertUserHandle finished")

	if err = ctx.BindJSON(&input); err != nil {
		log.Println(err)
		response.Errors = err.Error()
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	output, err := h.uc.UpdateUsers(ctx, input)
	if err != nil {
		log.Println(err)
		response.Errors = err.Error()
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	response.Message = "Created"
	response.Data = output.ID
	ctx.JSON(http.StatusCreated, response)
	return
}

func (h *Handler) ReadUsersHandle(ctx *gin.Context) {
	var (
		filter   entity.UserFilter
		response entity.Response
		err      error
	)
	log.Println("handle ReadUsersHandle started")
	defer log.Println("handle ReadUsersHandle finished")

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
	log.Println("handle DeleteUsersHandle started")
	defer log.Println("handle DeleteUsersHandle finished")

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
