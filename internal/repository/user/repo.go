package user

import (
	"context"
	"fmt"
	"tax-auth/internal/entity"
	"tax-auth/internal/repository"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

// https://github.com/Masterminds/squirrel

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

func (r *Repo) UpsertUser(ctx context.Context, user entity.User) (*entity.User, error) {
	var id int64
	const q = `INSERT INTO user ("name", "email", "password")
VALUES (@name, @email, @password)
ON CONFLICT (email, name)
    DO UPDATE SET email    = EXCLUDED.email,
                  name     = EXCLUDED.name,
                  password = EXCLUDED.password
RETURNING id;`
	args := pgx.NamedArgs{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	err := r.DB.QueryRow(ctx, q, args).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("unable to insert or update row: %v", err)
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
	).From("user")

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
		return nil, fmt.Errorf("unable to convert query to sql: %v", err)
	}

	rows, err := r.DB.Query(ctx, sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %v", err)
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
			return nil, fmt.Errorf("unable to scan row: %v", err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (r *Repo) DeleteUser(ctx context.Context, filter entity.UserFilter) error {
	q := r.builder.Delete("user")

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
		return fmt.Errorf("unable to convert query to sql: %v", err)
	}

	_, err = r.DB.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("unable to delete users: %v", err)
	}

	return nil
}
