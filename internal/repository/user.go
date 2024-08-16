package repository

import (
	"context"
	"fmt"
	"log"
	env "tax-auth/internal"
	"tax-auth/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

//https://github.com/Masterminds/squirrel

type UserRepo struct {
	Postgres
	builder squirrel.StatementBuilderType
}

type UserRepository interface {
	InsertUser(ctx context.Context, user entity.User) error
	ReadUsers(ctx context.Context, filter entity.Filter) ([]entity.User, error)
	UpdateUser(ctx context.Context, user entity.User, filter entity.Filter) error
	DeleteUser(ctx context.Context, filter entity.Filter) error
}

func NewUserRepo(ctx context.Context) UserRepo {
	url, err := env.GetDBUrlEnv()
	if err != nil {
		log.Fatal(err)
	}
	return UserRepo{
		Postgres: NewPG(ctx, url),
		builder:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (repo *UserRepo) InsertUser(ctx context.Context, user entity.User) error {
	q := `INSERT INTO users ("Name", "INN", "Email", "Password")
VALUES (@name, @inn, @email, @password);`
	args := pgx.NamedArgs{
		"name":     user.Name,
		"inn":      user.INN,
		"email":    user.Email,
		"password": user.Password,
	}
	tag, err := repo.db.Exec(ctx, q, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user '%s' with '%s' already exists", user.Name, user.Email)
	}
	return nil
}

func (repo *UserRepo) ReadUsers(ctx context.Context, filter entity.Filter) ([]entity.User, error) {
	var users []entity.User
	q := repo.builder.Select(
		"id",
		"name",
		"inn",
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

	rows, err := repo.db.Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.INN,
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

func (repo *UserRepo) UpdateUser(ctx context.Context, user entity.User, filter entity.Filter) error {

	return nil
}

func (repo *UserRepo) DeleteUser(ctx context.Context, filter entity.Filter) error {

	return nil
}
