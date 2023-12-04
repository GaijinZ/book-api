package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var dbPool *pgxpool.Pool

func getPostgresConnectionString() {
	var err error

	pgxAddress := os.Getenv("POSTGRES_BOOKS")

	dbPool, err = pgxpool.New(context.Background(), pgxAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err = dbPool.Ping(context.Background()); err != nil {
		log.Println(err)
	}

	log.Println("Connected!")
}

func GetConnection() *pgxpool.Pool {
	if dbPool == nil {
		getPostgresConnectionString()
	}

	return dbPool
}
