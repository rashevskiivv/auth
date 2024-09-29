package auth

import (
	"context"
	"errors"
	"log"
	"net/mail"
	env "tax-auth/internal"
	"tax-auth/internal/entity"
	"tax-auth/internal/repository/auth"
	repositoryUser "tax-auth/internal/repository/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repoToken auth.Repository
	repoUser  repositoryUser.Repository
}

func NewAuthUseCase(repo auth.Repository, repoUser repositoryUser.Repository) *UseCase {
	return &UseCase{
		repoToken: repo,
		repoUser:  repoUser,
	}
}

func (uc *UseCase) RegisterUser(ctx context.Context, input entity.RegisterInput) (*entity.RegisterOutput, error) {
	emailOk, err := validateEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if !emailOk {
		return nil, errors.New("email is required")
	}
	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	users, err := uc.repoUser.ReadUsers(ctx, entity.UserFilter{
		Email: []string{input.Email},
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(users) != 0 {
		return nil, errors.New("user with this email already registered")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return nil, err
	}
	hashedPassword := string(bytes)
	input.User.Password = hashedPassword

	user, err := uc.repoUser.InsertUser(ctx, input.User) //todo add dbtx
	if err != nil {
		return nil, err
	}
	if user.ID == nil || *user.ID == 0 {
		return nil, errors.New("user has no id")
	}

	secretKey, err := env.GetJWTSecretKey()
	if err != nil {
		return nil, err
	}
	payload := jwt.MapClaims{
		"sub": input.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	log.Println(t)

	insertTokenInput := entity.Token{
		Token:  t,
		UserID: *user.ID,
	}
	err = uc.repoToken.InsertToken(ctx, insertTokenInput)
	if err != nil {
		return nil, err
	}

	response := entity.RegisterOutput{
		Token: entity.Token{
			Token:  t,
			UserID: *user.ID,
		},
	}
	return &response, nil
}

func (uc *UseCase) AuthenticateUser(ctx context.Context, input entity.AuthenticateInput) (*entity.AuthenticateOutput, error) {
	err := bcrypt.CompareHashAndPassword([]byte(input.Hash), []byte(input.Password))
	return nil, err
}

func validateEmail(email string) (bool, error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return false, err
	}
	return address.Address != "", nil
}
