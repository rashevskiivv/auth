package auth

import (
	"log"
	"net/http"
	"tax-auth/internal/entity"
	auth "tax-auth/internal/usecase/auth"

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
	request := entity.User{}
	err = ctx.ShouldBind(&request)
	if err != nil {
		log.Println(err)                                             // todo test log.Fatal and Println
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //todo change to response
		return
	}

	registerInput := entity.RegisterInput{User: request}
	token, err := h.uc.RegisterUser(ctx, registerInput)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) //todo change to response struct
		return
	}
	log.Println(token)

	log.Println("registered")
	ctx.JSON(http.StatusCreated, entity.Response{
		Data:    token,
		Message: "User registered. Access token in data.",
		Errors:  "",
	})
	return
	/* todo
	1. get login and password
	2. validate them
	3. find same login in db
	4. hash last one
	5. save to the db user
	6. generate jwt token
	7. save to the db jwt token
	8. return jwt token
	*/
}

func (h *Handler) AuthenticateUserHandle(ctx *gin.Context) {
	log.Println("logged in")
	ctx.JSON(http.StatusCreated, "logged in")
}
