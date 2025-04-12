package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/rashevskiivv/auth/internal/entity"
	"github.com/rashevskiivv/auth/internal/repository/auth"
	repositoryUser "github.com/rashevskiivv/auth/internal/repository/user"
	"github.com/rashevskiivv/auth/internal/usecase"

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
	emailOk, err := usecase.ValidateEmail(input.Email)
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

	hashedPassword, err := usecase.GenerateHashedPassword(input.Password)
	if err != nil {
		return nil, err
	}
	input.User.Password = *hashedPassword

	user, err := uc.repoUser.UpsertUser(ctx, input.User) //todo add dbtx
	if err != nil {
		return nil, err
	}
	if user.ID == nil || *user.ID == 0 {
		return nil, errors.New("user has no id")
	}

	t, err := usecase.GetJWTToken(input.Email)
	if err != nil {
		return nil, err
	}

	insertTokenInput := entity.Token{
		Token:  *t,
		UserID: *user.ID,
	}
	err = uc.repoToken.InsertToken(ctx, insertTokenInput)
	if err != nil {
		return nil, err
	}

	response := entity.RegisterOutput{
		Token: entity.Token{
			Token:  *t,
			UserID: *user.ID,
		},
	}
	return &response, nil
}

func (uc *UseCase) AuthenticateUser(ctx context.Context, input entity.AuthenticateInput) (*entity.AuthenticateOutput, error) {
	emailOk, err := usecase.ValidateEmail(input.Email)
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
	if len(users) == 0 {
		return nil, errors.New("no user with this email")
	}
	user := users[0]
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, err
	}

	t, err := usecase.GetJWTToken(input.Email)
	if err != nil {
		return nil, err
	}

	insertTokenInput := entity.Token{
		Token:  *t,
		UserID: *user.ID,
	}
	err = uc.repoToken.InsertToken(ctx, insertTokenInput)
	if err != nil {
		return nil, err
	}

	response := entity.AuthenticateOutput{
		Token: entity.Token{
			Token:  *t,
			UserID: *user.ID,
		},
	}
	return &response, nil
}

func (uc *UseCase) CheckToken(ctx context.Context, input entity.CheckTokenInput) (entity.CheckTokenOutput, error) {
	if input.UserID == "" {
		return entity.CheckTokenOutput{}, fmt.Errorf("id is empty")
	}
	_, err := strconv.ParseInt(input.UserID, 10, 64)
	if err != nil {
		log.Println(err)
		return entity.CheckTokenOutput{}, err
	}

	filter := entity.TokenFilter{UserID: []string{input.UserID}}
	tokens, err := uc.repoToken.ReadTokens(ctx, filter)
	if err != nil {
		return entity.CheckTokenOutput{}, err
	}
	if len(tokens) == 0 {
		return entity.CheckTokenOutput{}, errors.New("no token with specified user_id")
	}

	return entity.CheckTokenOutput{Token: tokens[0]}, nil
}
