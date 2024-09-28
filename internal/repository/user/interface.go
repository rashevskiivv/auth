package user

import (
	"context"
	"tax-auth/internal/entity"
)

type Repository interface {
	InsertUser(ctx context.Context, user entity.User) error
	ReadUsers(ctx context.Context, filter entity.Filter) ([]entity.User, error)
	UpdateUser(ctx context.Context, user entity.User, filter entity.Filter) error
	DeleteUser(ctx context.Context, filter entity.Filter) error
}
