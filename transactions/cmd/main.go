package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"library/pkg/config"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/pkg/rabbitMQ/rabbitMQ"
	"library/pkg/redis"
	"library/pkg/utils"
	"library/transactions/handler"
	"library/transactions/repository"
	"library/transactions/server"
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

	redisClient, err := redis.NewRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	rmq, err := rabbitMQ.NewConn(cfg)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ instance: %v", err)
	}

	transactionsRepository := repository.NewTransactionRepository(ctx, *db, redisClient, rmq)
	transactionsHandler := handler.NewTransactionHandler(ctx, transactionsRepository)

	router := server.NewRouter(transactionsHandler)

	go router.Run(":" + cfg.TransactionsServerPort)

	defer func() {
		cancel()
		db.Close()
		redisClient.Close()
		rmq.Close()
	}()

	select {
	case sig := <-interrupt:
		log.Infof("Received signal: %v\n", sig)
	case <-ctx.Done():
		log.Infof("Context done")
	}
}
