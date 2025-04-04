package auth

import (
	"context"
	"tax-auth/internal/entity"
)

type UseCaseI interface {
	RegisterUser(ctx context.Context, input entity.RegisterInput) (*entity.RegisterOutput, error)
	AuthenticateUser(ctx context.Context, input entity.AuthenticateInput) (*entity.AuthenticateOutput, error)
	CheckToken(ctx context.Context, input entity.CheckTokenInput) (entity.CheckTokenOutput, error)
}
