package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"

	"library/pkg/config"
	"library/pkg/middleware"
	"library/pkg/postgres"
	"library/pkg/tracing"
	"library/pkg/utils"
	"library/users/handler"
	"library/users/repository"

	"log"
	"time"
)

func Run(ctx *context.Context, cfg config.GlobalEnv) {
	if err := envconfig.Process("bookapi", &cfg); err != nil {
		log.Fatal(err.Error())
	}

	log := utils.GetLogger(*ctx)

	configDB := postgres.DBConfig{
		DriverName:      "postgres",
		DataSourceName:  cfg.PostgresBooks,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}

	db, err := postgres.NewDB(*ctx, configDB)
	if err != nil {
		log.Errorf("Failed to configure db connection: %v", err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository(*ctx, *db)
	authRepository := repository.NewAuthRepository(*ctx, *db)
	authUser := handler.NewUserAuth(*ctx, authRepository)
	handlerUser := handler.NewUserHandler(*ctx, userRepository)

	router := gin.Default()

	router.POST("/v1/users", handlerUser.AddUser)
	router.POST("/v1/users/login", authUser.Login)
	router.POST("/v1/users/logout", authUser.Logout)

	v1 := router.Group("/v1/users")

	v1.PUT("/:user_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		handlerUser.UpdateUser,
	)
	v1.GET("/:user_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		handlerUser.GetUser,
	)
	v1.GET("",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		handlerUser.GetAllUsers,
	)
	v1.DELETE("/:user_id/:delete_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		middleware.GetDeleteParam,
		handlerUser.DeleteUser,
	)

	router.Run(":" + cfg.UsersServerPort)
}
