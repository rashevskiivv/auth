package auth

import (
	"context"
	"tax-auth/internal/entity"
)

type Repository interface {
	ReadTokens(ctx context.Context, filter entity.TokenFilter) ([]entity.Token, error)
	InsertToken(ctx context.Context, token entity.Token) error
}
