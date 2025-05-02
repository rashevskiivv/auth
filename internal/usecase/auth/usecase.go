package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"

	env "github.com/rashevskiivv/auth/internal"
	"github.com/rashevskiivv/auth/internal/client"
	"github.com/rashevskiivv/auth/internal/entity"
	"github.com/rashevskiivv/auth/internal/repository/auth"
	repositoryUser "github.com/rashevskiivv/auth/internal/repository/user"
	"github.com/rashevskiivv/auth/internal/usecase"

	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	client    *client.Client
	repoToken auth.Repository
	repoUser  repositoryUser.Repository
}

func NewAuthUseCase(repo auth.Repository, repoUser repositoryUser.Repository) *UseCase {
	return &UseCase{
		client:    client.NewClient(),
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

	user, err := uc.repoUser.UpsertUser(ctx, input.User) // todo add dbtx
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

	input.Token = *t
	input.RequestUtils.ID = strconv.FormatInt(*user.ID, 10)
	err = uc.makeRequests(input)
	if err != nil {
		return nil, err
	}

	response := entity.RegisterOutput{
		Token: entity.Token{
			ID:    user.ID,
			Token: *t,
		},
	}
	return &response, nil
}

func (uc *UseCase) makeRequests(input entity.RegisterInput) error {
	var (
		err      error
		appURL   string
		apiReq   *client.Request
		recomReq *client.Request
	)
	switch input.WhichRequest {
	case "":
		appURL, err = env.GetAPIAppURL()
		if err != nil {
			log.Println(err)
			return err
		}
		apiReq, err = buildReq(input, appURL)
		if err != nil {
			return err
		}

		appURL, err = env.GetRecommendationAppURL()
		if err != nil {
			log.Println(err)
			return err
		}
		recomReq, err = buildReq(input, appURL)
		if err != nil {
			return err
		}
	case entity.AppRecommendations:
		appURL, err = env.GetAPIAppURL()
		if err != nil {
			log.Println(err)
			return err
		}
		apiReq, err = buildReq(input, appURL)
		if err != nil {
			return err
		}
	case entity.AppAPI:
		appURL, err = env.GetRecommendationAppURL()
		if err != nil {
			log.Println(err)
			return err
		}
		recomReq, err = buildReq(input, appURL)
		if err != nil {
			return err
		}
	default:
		return nil
	}
	err = uc.sendBothReqs(apiReq, recomReq)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) sendBothReqs(apiReq *client.Request, recomReq *client.Request) error {
	var apiID, recomID int64 = -1, -2
	for _, req := range []*client.Request{apiReq, recomReq} {
		if req == nil {
			continue
		}

		resp, err := uc.client.Do(req)
		if err != nil {
			log.Println(err)
			return err
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(resp.Body)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return err
		}

		response := entity.Response{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Println(err)
			return err
		}
		idf := response.Data.(float64)
		if apiID < 0 {
			apiID = int64(math.Round(idf))
		} else {
			recomID = int64(math.Round(idf))
		}
	}
	if apiID != -1 && recomID != -2 && apiID != recomID {
		return fmt.Errorf("ids are not the same")
	}
	return nil
}

func buildReq(input entity.RegisterInput, appURL string) (*client.Request, error) {
	out, err := json.Marshal(input.User)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req := client.NewRequest(http.MethodPost, appURL+entity.PathUsers, bytes.NewBuffer(out))
	if req == nil {
		log.Printf("req with url %v is nil", appURL)
		return nil, fmt.Errorf("req with url %v is nil", appURL)
	}

	headers := make(map[string]string, 3)
	headers["id"] = input.RequestUtils.ID
	headers["token"] = input.Token
	headers["Origin"] = entity.AppAuth
	req.AddAuthHeaders(headers)

	return req, nil
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
