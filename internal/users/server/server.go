package server

import (
	"library/database/postgres"
	"library/internal/users/handler"
	"library/internal/users/middleware"
	"library/internal/users/repository"

	"github.com/gin-gonic/gin"
)

func Run(port string) {
	dbPool := postgres.GetPostgresConnectionString()
	defer dbPool.Close()

	userRepository := repository.NewUserRepository(dbPool)
	authRepository := repository.NewAuthRepository(dbPool)
	authUser := handler.NewUserAuth(authRepository)
	handlerUser := handler.NewUserHandler(userRepository)

	router := gin.Default()

	router.POST("/v1/users/login", authUser.Login)
	router.POST("/v1/users/logout", authUser.Logout)

	v1 := router.Group("/v1")
	v1.Use(middleware.IsAuthorized())

	v1.POST("/user", handlerUser.AddUser)
	v1.PUT("/user/:user_id", handlerUser.UpdateUser)
	v1.GET("/user/:user_id", handlerUser.GetUser)
	v1.GET("/user", handlerUser.GetAllUsers)
	v1.DELETE("/user/:user_id", handlerUser.DeleteUser)

	router.Run(port)
}
