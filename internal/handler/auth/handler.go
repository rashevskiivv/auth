package auth

import (
	"log"
	"net/http"
	"tax-auth/internal/repository/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo auth.Repository
}

func NewAuthHandler(repo auth.Repository) Handler {
	return Handler{
		repo: repo,
	}
}

func (h *Handler) RegisterUserHandle(ctx *gin.Context) {
	log.Println("registered")
	ctx.JSON(http.StatusCreated, "registered")
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func comparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
