package auth

import (
	"context"
	"tax-auth/internal/entity"
)

type Repository interface {
	ReadToken(ctx context.Context, filter entity.Filter) error
	InsertToken(ctx context.Context, token string) error
}
