package auth

import (
	"context"
	"tax-auth/internal/entity"
)

type Repository interface {
	ReadToken(ctx context.Context, filter entity.UserFilter) error
	InsertToken(ctx context.Context, token entity.Token) error
}
