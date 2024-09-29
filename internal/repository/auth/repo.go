package auth

import (
	"context"
	"tax-auth/internal/entity"
	"tax-auth/internal/repository"

	"github.com/Masterminds/squirrel"
)

type Repo struct {
	repository.Postgres
	builder squirrel.StatementBuilderType
}

func NewAuthRepo(pg repository.Postgres) *Repo {
	return &Repo{
		Postgres: pg,
		builder:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (repo *Repo) ReadToken(ctx context.Context, filter entity.UserFilter) error {
	return nil
}

func (repo *Repo) InsertToken(ctx context.Context, token entity.Token) error {
	return nil
}
