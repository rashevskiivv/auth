package user

import (
	"context"

	"github.com/rashevskiivv/auth/internal/entity"
)

type Repository interface {
	UpsertUser(ctx context.Context, user entity.User) (*entity.User, error)
	ReadUsers(ctx context.Context, filter entity.UserFilter) ([]entity.User, error)
	DeleteUser(ctx context.Context, filter entity.UserFilter) error
}
