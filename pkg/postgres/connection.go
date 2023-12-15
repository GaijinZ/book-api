package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool" // better to try use standard sql driver, in thus case you can migrate from one database to another
	"github.com/kelseyhightower/envconfig"
	"library/pkg/config"
	"library/pkg/logger"
	"log"
	"os"
)

var dbPool *pgxpool.Pool // don't use a singleton pattern for database connection

// Make a DB struck, with configured values
// add ping, close, getConnection funcs

func getPostgresConnectionString() {
	var cfg config.GlobalEnv
	var err error

	if err := envconfig.Process("bookapi", &cfg); err != nil { // don't read config on pkg level, read config in main() and pass in as a value to every package
		log.Fatal(err.Error())
	}

	log := logger.NewLogger() //don't create a new logger, reuse already created

	dbPool, err = pgxpool.New(context.Background(), cfg.PostgresBooks)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err = dbPool.Ping(context.Background()); err != nil { // pass context as a value
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
