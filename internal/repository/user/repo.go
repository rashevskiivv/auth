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

func (r *Repo) InsertUser(ctx context.Context, user entity.User) (*entity.User, error) {
	var id int64
	q := `INSERT INTO users ("name", "email", "password")
VALUES (@name, @email, @password)
RETURNING id;`
	args := pgx.NamedArgs{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	err := r.DB.QueryRow(ctx, q, args).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("unable to insert row: %w", err)
	}
	return &entity.User{ID: &id}, nil
}

func (r *Repo) ReadUsers(ctx context.Context, filter entity.UserFilter) ([]entity.User, error) {
	var users []entity.User
	q := r.builder.Select(
		"id",
		"name",
		"email",
		"password",
	).From("users")

	// Where
	if len(filter.ID) > 0 {
		q = q.Where(squirrel.Eq{"id": filter.ID})
	}
	if len(filter.Email) > 0 {
		q = q.Where(squirrel.Eq{"email": filter.Email})
	}
	if len(filter.Name) > 0 {
		q = q.Where(squirrel.Eq{"name": filter.Name})
	}

	// Limit
	if filter.Limit != 0 {
		q = q.Limit(uint64(filter.Limit))
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("unable to convert query to sql: %w", err)
	}

	rows, err := r.DB.Query(ctx, sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}

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

func (r *Repo) UpdateUser(ctx context.Context, user entity.User, filter entity.UserFilter) error {

	return nil
}

func (r *Repo) DeleteUser(ctx context.Context, filter entity.UserFilter) error {

	return nil
}
