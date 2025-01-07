package repository

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

//https://donchev.is/post/working-with-postgresql-in-go-using-pgx/

type Postgres struct {
	DB *pgxpool.Pool
}

var (
	pgInstance Postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) (*Postgres, error) {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Panicf("unable to create connection pool: %s", err)
		}

		pgInstance = Postgres{DB: db}
	})

	err := pgInstance.Ping(ctx)
	if err != nil {
		log.Panicf("unable to ping db: %s", err)
		return nil, err
	}
	return &pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() error {
	pg.DB.Close()
	return nil
}
