package user

import (
	"context"
	"tax-auth/internal/entity"
)

type Repository interface {
	InsertUser(ctx context.Context, user entity.User) (*entity.User, error)
	ReadUsers(ctx context.Context, filter entity.UserFilter) ([]entity.User, error)
	UpdateUser(ctx context.Context, user entity.User, filter entity.UserFilter) error
	DeleteUser(ctx context.Context, filter entity.UserFilter) error
}
