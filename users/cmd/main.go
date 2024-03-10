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

//	@Title		UserAPI Service
//	@Version	1.0
//	@Host		localhost:5000
//
// BasePath /v1/users
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

	userRepository := repository.NewUserRepository(ctx, *db, redisClient, rmq)
	authRepository := repository.NewAuthRepository(ctx, *db)
	authUser := handler.NewUserAuth(ctx, authRepository)
	handlerUser := handler.NewUserHandler(ctx, userRepository)

	router := server.NewRouter(authUser, handlerUser)

	go router.Run(":" + cfg.UsersServerPort)

	defer func() {
		cancel()
		db.Close()
		redisClient.Close()
		rmq.Close()
	}()

	select {
	case <-interrupt:
		log.Infof("Received a shutdown signal...")
		close(interrupt)
	case <-ctx.Done():
		log.Infof("Context done")
		close(interrupt)
	}
}
