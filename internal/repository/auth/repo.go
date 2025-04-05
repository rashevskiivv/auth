package auth

import (
	"context"
	"fmt"

	"github.com/rashevskiivv/auth/internal/entity"
	"github.com/rashevskiivv/auth/internal/repository"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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

func (r *Repo) ReadTokens(ctx context.Context, filter entity.TokenFilter) ([]entity.Token, error) {
	var tokens []entity.Token
	q := r.builder.Select(
		"id",
		"token",
		"user_id",
	).From("token")

	// Where
	if len(filter.ID) > 0 {
		q = q.Where(squirrel.Eq{"id": filter.ID})
	}
	if len(filter.Token) > 0 {
		q = q.Where(squirrel.Eq{"token": filter.Token})
	}
	if len(filter.UserID) > 0 {
		q = q.Where(squirrel.Eq{"user_id": filter.UserID})
	}

	// Limit
	if filter.Limit != 0 {
		q = q.Limit(uint64(filter.Limit))
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to convert query to sql: %v", err)
	}

	rows, err := r.DB.Query(ctx, sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to query tokens: %v", err)
	}

	for rows.Next() {
		token := entity.Token{}
		err = rows.Scan(
			&token.ID,
			&token.Token,
			&token.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		tokens = append(tokens, token)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return tokens, nil
}

func (r *Repo) InsertToken(ctx context.Context, token entity.Token) error {
	q := `INSERT INTO token ("token", "user_id")
VALUES (@token, @user_id)
ON CONFLICT (user_id)
    DO UPDATE SET token   = EXCLUDED.token,
                  user_id = EXCLUDED.user_id;`
	args := pgx.NamedArgs{
		"token":   token.Token,
		"user_id": token.UserID,
	}
	_, err := r.DB.Query(ctx, q, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %v", err)
	}
	return nil
}
