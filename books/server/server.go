package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/books/handler"
	"library/books/repository"
	"library/pkg/middleware"
	"library/pkg/postgres"
	"library/pkg/tracing"
)

func Run(ctx *context.Context, port string) {
	dbPool := postgres.GetConnection()
	defer dbPool.Close()

	bookRepository := repository.NewBookRepository(*ctx, dbPool)
	handlerBook := handler.NewBookHandler(*ctx, bookRepository)

	router := gin.Default()

	v1 := router.Group("/v1/books")

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
		handlerBook.UpdateBook,
	)
	v1.GET("/:user_id/books/:book_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
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
		handlerBook.DeleteBook,
	)

	router.Run(port)
}
