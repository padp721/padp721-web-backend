package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
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
		log.Fatal(err.Error())
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err.Error())
	}

	var pgVersion string
	sql := "SELECT version();"
	err = dbPool.QueryRow(context.Background(), sql).Scan(&pgVersion)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(pgVersion)

	return dbPool
}
