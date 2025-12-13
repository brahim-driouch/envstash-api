package config

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func ConnectDB() error {
	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		return errors.New("DATABASE_URL is not set")
	}
	pool, err := pgxpool.New(context.TODO(), dbUrl)
	if err != nil {
		return err
	}

	Pool = pool
	return nil
}
