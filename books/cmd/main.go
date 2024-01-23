package main

import (
	"context"

	"github.com/kelseyhightower/envconfig"

	"library/books/handler"
	"library/books/repository"
	"library/books/server"
	"library/pkg/config"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/pkg/utils"

	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var cfg config.GlobalEnv
	var ctx context.Context

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))
	log := utils.GetLogger(ctx)

	if err := envconfig.Process("bookapi", &cfg); err != nil {
		log.Fatalf(err.Error())
	}

	configDB := postgres.DBConfig{
		DriverName:      "postgres",
		DataSourceName:  cfg.PostgresBooks,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}

	db, err := postgres.NewDB(ctx, configDB)
	if err != nil {
		log.Errorf("Failed to configure db connection: %v", err)
	}
	defer db.Close()

	bookRepository := repository.NewBookRepository(ctx, *db)
	handlerBook := handler.NewBookHandler(ctx, bookRepository)

	router := server.NewRouter(handlerBook)

	go router.Run(":" + cfg.BooksServerPort)

	select {
	case sig := <-interrupt:
		log.Infof("Received signal: %v\n", sig)
	case <-ctx.Done():
		log.Infof("Context done")
	}
}
