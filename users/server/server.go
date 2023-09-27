package server

import (
	"github.com/gin-gonic/gin"
	middleware "library/pkg/middleware"
	handler2 "library/users/handler"
	"library/users/postgres"
	repository2 "library/users/repository"
)

func Run(port string) {
	dbPool := postgres.GetPostgresConnectionString()
	defer dbPool.Close()

	userRepository := repository2.NewUserRepository(dbPool)
	authRepository := repository2.NewAuthRepository(dbPool)
	authUser := handler2.NewUserAuth(authRepository)
	handlerUser := handler2.NewUserHandler(userRepository)

	router := gin.Default()

	router.POST("/v1/users", handlerUser.AddUser)
	router.POST("/v1/users/login", authUser.Login)
	router.POST("/v1/users/logout", authUser.Logout)

	v1 := router.Group("/v1")
	v1.Use(middleware.IsAuthorized())

	users := v1.Group("/users")

	users.PUT("/:user_id", handlerUser.UpdateUser)
	users.GET("/:user_id", handlerUser.GetUser)
	users.GET("", handlerUser.GetAllUsers)
	users.DELETE("/:user_id/:delete_id", handlerUser.DeleteUser)

	router.Run(port)
}
