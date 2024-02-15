package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"library/books/handler"

	"library/pkg/middleware"
	"library/pkg/tracing"
)

func NewRouter(handlerBook handler.BookerHandler) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1/books")

	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.POST("/:user_id/books",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		handlerBook.AddBook,
	)
	v1.PUT("/:user_id/books/:book_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		middleware.GetBookParam,
		handlerBook.UpdateBook,
	)
	v1.GET("/:user_id/books/:book_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		middleware.GetBookParam,
		handlerBook.GetBook,
	)
	v1.GET("/:user_id/books",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		handlerBook.GetAllBooks,
	)
	v1.DELETE("/:user_id/books/:book_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		middleware.GetBookParam,
		handlerBook.DeleteBook,
	)

	return router
}
