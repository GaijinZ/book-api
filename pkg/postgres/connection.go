package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"library/pkg/config"
	"library/pkg/logger"
	"log"
	"os"
)

var dbPool *pgxpool.Pool

func getPostgresConnectionString() {
	var cfg config.GlobalEnv
	var err error

	if err := envconfig.Process("bookapi", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	log := logger.NewLogger()

	dbPool, err = pgxpool.New(context.Background(), cfg.PostgresBooks)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err = dbPool.Ping(context.Background()); err != nil {
		log.Errorf("Ping failed: %v", err)
	}

	log.Infof("Postgres connected!")
}

func GetConnection() *pgxpool.Pool {
	if dbPool == nil {
		getPostgresConnectionString()
	}

	return dbPool
}
