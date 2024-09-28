package user

import (
	"context"
	"fmt"
	"tax-auth/internal/entity"
	"tax-auth/internal/repository"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

//https://github.com/Masterminds/squirrel

type Repo struct {
	repository.Postgres
	builder squirrel.StatementBuilderType
}

func NewUserRepo(pg repository.Postgres) *Repo {
	return &Repo{
		Postgres: pg,
		builder:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (repo *Repo) InsertUser(ctx context.Context, user entity.User) error {
	q := `INSERT INTO users ("Name", "Email", "Password")
VALUES (@name, @email, @password);`
	args := pgx.NamedArgs{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	tag, err := repo.DB.Exec(ctx, q, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user '%s' with '%s' already exists", user.Name, user.Email)
	}
	return nil
}

func (repo *Repo) ReadUsers(ctx context.Context, filter entity.Filter) ([]entity.User, error) {
	var users []entity.User
	q := repo.builder.Select(
		"id",
		"name",
		"email",
		"password",
	).From(`users`)

	if len(filter.Conditions) > 0 { //todo len(nil) == ?
		for key, values := range filter.Conditions {
			q = q.Where(squirrel.Eq{key: values})
		}
	}

	if filter.Limit != 0 {
		q = q.Limit(uint64(filter.Limit))
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to convert query to sql: %w", err)
	}

	rows, err := repo.DB.Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
		)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (repo *Repo) UpdateUser(ctx context.Context, user entity.User, filter entity.Filter) error {

	return nil
}

func (repo *Repo) DeleteUser(ctx context.Context, filter entity.Filter) error {

	return nil
}
