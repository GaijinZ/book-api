package server

import (
	"github.com/gin-gonic/gin"
	"library/pkg/middleware"
	"library/pkg/tracing"
	"library/users/handler"
)

func NewRouter(authUser handler.UserAuther, handlerUser handler.Userer) *gin.Engine {
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

	return router
}
