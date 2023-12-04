package server

import (
	"github.com/gin-gonic/gin"
	middleware "library/pkg/middleware"
	"library/pkg/postgres"
	"library/users/handler"
	repository2 "library/users/repository"
)

func Run(port string) {
	dbPool := postgres.GetConnection()
	defer dbPool.Close()

	userRepository := repository2.NewUserRepository(dbPool)
	authRepository := repository2.NewAuthRepository(dbPool)
	authUser := handler.NewUserAuth(authRepository)
	handlerUser := handler.NewUserHandler(userRepository)

	router := gin.Default()

	router.POST("/v1/users", handlerUser.AddUser)
	router.POST("/v1/users/login", authUser.Login)
	router.POST("/v1/users/logout", authUser.Logout)

	v1 := router.Group("/v1/users")

	v1.PUT("/:user_id", middleware.IsAuthorized, handlerUser.UpdateUser)
	v1.GET("/:user_id", middleware.IsAuthorized, handlerUser.GetUser)
	v1.GET("", middleware.IsAuthorized, handlerUser.GetAllUsers)
	v1.DELETE("/:user_id/:delete_id", middleware.IsAuthorized, handlerUser.DeleteUser)

	router.Run(port)
}
