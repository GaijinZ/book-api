package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/pkg/middleware"
	"library/pkg/postgres"
	"library/pkg/tracing"
	"library/users/handler"
	repository2 "library/users/repository"
)

func Run(ctx context.Context, port string) {
	dbPool := postgres.GetConnection()
	defer dbPool.Close()

	userRepository := repository2.NewUserRepository(ctx, dbPool)
	authRepository := repository2.NewAuthRepository(ctx, dbPool)
	authUser := handler.NewUserAuth(ctx, authRepository)
	handlerUser := handler.NewUserHandler(ctx, userRepository)

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
		handlerUser.DeleteUser,
	)

	router.Run(port)
}
