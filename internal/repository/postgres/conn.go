package postgres
// db connection placeholder

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	dsn := "postgres://app:app@localhost:5432/mydb?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}