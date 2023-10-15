package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

func GetPostgresConnectionString() *pgxpool.Pool {
	pgxAddress := os.Getenv("POSTGRES_SHOPS")

	dbPool, err := pgxpool.New(context.Background(), pgxAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Println(err)
	}

	log.Println("Connected!")

	return dbPool
}
