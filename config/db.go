package config

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool = nil

func ConnectDB() error {
	if Pool != nil {
		return nil // Already connected
	}
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return errors.New("database url environment variable is not set")
	}
	pool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return err
	}
	if pool == nil {
		return errors.New("failed to create database connection pool")
	}
	Pool = pool
	return nil
}
