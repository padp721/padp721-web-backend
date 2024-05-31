package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePool() (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v pool_max_conns=5",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	var pgVersion string
	sql := "SELECT version();"
	err = dbPool.QueryRow(context.Background(), sql).Scan(&pgVersion)
	if err != nil {
		return nil, err
	}
	log.Println(pgVersion)

	return dbPool, nil
}
