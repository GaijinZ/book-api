package main

import (
	"context"
	"library/pkg/config"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/pkg/rabbitMQ/rabbitMQ"
	"library/pkg/redis"
	"library/pkg/utils"
	"library/users/handler"
	"library/users/repository"
	"library/users/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
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

	redisClient, err := redis.NewRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	rmq, err := rabbitMQ.NewConn()
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ instance: %v", err)
	}
	defer rmq.Close()

	userRepository := repository.NewUserRepository(ctx, *db, redisClient, rmq)
	authRepository := repository.NewAuthRepository(ctx, *db)
	authUser := handler.NewUserAuth(ctx, authRepository)
	handlerUser := handler.NewUserHandler(ctx, userRepository)

	router := server.NewRouter(authUser, handlerUser)

	go router.Run(":" + cfg.UsersServerPort)

	select {
	case <-interrupt:
		log.Infof("Received a shutdown signal...")
		close(interrupt)
	case <-ctx.Done():
		log.Infof("Context done")
		close(interrupt)
	}
}
