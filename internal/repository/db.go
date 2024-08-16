package repository

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

//https://donchev.is/post/working-with-postgresql-in-go-using-pgx/

type Postgres struct {
	db *pgxpool.Pool
}

var (
	pgInstance Postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) Postgres {
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			log.Panicf("unable to create connection pool: %s", err)
		}

		pgInstance = Postgres{db}
	})

	err := pgInstance.Ping(ctx)
	if err != nil {
		log.Panicf("unable to ping db: %s", err)
	}
	return pgInstance
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.db.Close()
}
