package server

import (
	"github.com/gin-gonic/gin"
	"library/books/handler"
	"library/books/postgres"
	"library/books/repository"
	middleware "library/pkg/middleware"
)

func Run(port string) {
	dbPool := postgres.GetPostgresConnectionString()
	defer dbPool.Close()

	bookRepository := repository.NewBookRepository(dbPool)
	handlerBook := handler.NewBookHandler(bookRepository)

	router := gin.Default()

	v1 := router.Group("/v1/books")

	v1.POST("/:user_id/books", middleware.IsAuthorized, handlerBook.AddBook)
	v1.PUT("/:user_id/books/:book_id", middleware.IsAuthorized, handlerBook.UpdateBook)
	v1.GET("/:user_id/books/:book_id", middleware.IsAuthorized, handlerBook.GetBook)
	v1.GET("/:user_id/books", middleware.IsAuthorized, handlerBook.GetAllBooks)
	v1.DELETE("/:user_id/books/:book_id", middleware.IsAuthorized, handlerBook.DeleteBook)

	router.Run(port)
}
