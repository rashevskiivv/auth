package user

import (
	"context"
	"tax-auth/internal/entity"
)

type UseCaseI interface {
	GetUsers(ctx context.Context, input entity.GetUsersInput) (*entity.GetUsersOutput, error)
	UpdateUsers(ctx context.Context, input entity.UpdateUsersInput) error
	DeleteUsers(ctx context.Context, input entity.DeleteUsersInput) error
}
