package server

import (
	"userapi/internal/users/db/postgres"
	"userapi/internal/users/handler"
	"userapi/internal/users/repository"

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
	v1 := router.Group("/v1")
	v1.POST("/users/login", authUser.Login)
	v1.POST("/users", handlerUser.AddUser)
	v1.PUT("/users/:id", handlerUser.UpdateUser)
	v1.GET("/users/:id", handlerUser.GetUser)
	v1.GET("/users", handlerUser.GetAllUsers)
	v1.DELETE("/users/:id", handlerUser.DeleteUser)

	router.Run(port)
}
