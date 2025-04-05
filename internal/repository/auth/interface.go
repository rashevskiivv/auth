package auth

import (
	"context"

	"github.com/rashevskiivv/auth/internal/entity"
)

type Repository interface {
	ReadTokens(ctx context.Context, filter entity.TokenFilter) ([]entity.Token, error)
	InsertToken(ctx context.Context, token entity.Token) error
}
