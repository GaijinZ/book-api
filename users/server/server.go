package server

import (
	"github.com/gin-gonic/gin"
	"library/pkg/middleware"
	"library/pkg/tracing"
	"library/users/handler"
)

func NewRouter(authUser handler.UserAuther, handlerUser handler.Userer) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1/users")

	v1.POST("", handlerUser.AddUser)
	v1.GET("/activate", handlerUser.ActivateAccount)
	v1.POST("/login", authUser.Login)
	v1.POST("/logout", authUser.Logout)

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
