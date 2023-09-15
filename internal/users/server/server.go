package server

import (
	"github.com/gin-gonic/gin"
	"library/internal/users/handler"
	"library/internal/users/postgres"
	"library/internal/users/repository"
	middleware "library/utils/middleware"
)

func Run(port string) {
	dbPool := postgres.GetPostgresConnectionString()
	defer dbPool.Close()

	userRepository := repository.NewUserRepository(dbPool)
	authRepository := repository.NewAuthRepository(dbPool)
	authUser := handler.NewUserAuth(authRepository)
	handlerUser := handler.NewUserHandler(userRepository)

	router := gin.Default()

	router.POST("/v1/users", handlerUser.AddUser)
	router.POST("/v1/users/login", authUser.Login)
	router.POST("/v1/users/logout", authUser.Logout)

	v1 := router.Group("/v1")
	v1.Use(middleware.IsAuthorized())

	v1.PUT("/users/:user_id", handlerUser.UpdateUser)
	v1.GET("/users/:user_id", handlerUser.GetUser)
	v1.GET("/users", handlerUser.GetAllUsers)
	v1.DELETE("/users/:user_id/:delete_id", handlerUser.DeleteUser)

	router.Run(port)
}
