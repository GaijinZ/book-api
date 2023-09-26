package server

import (
	"library/internal/books/handler"
	"library/internal/books/postgres"
	"library/internal/books/repository"
	middleware "library/utils/middleware"

	"github.com/gin-gonic/gin"
)

func Run(port string) {
	dbPool := postgres.GetPostgresConnectionString()
	defer dbPool.Close()

	bookRepository := repository.NewBookRepository(dbPool)
	handlerBook := handler.NewBookHandler(bookRepository)

	router := gin.Default()

	v1 := router.Group("/v1")
	v1.Use(middleware.IsAuthorized())

	books := v1.Group("/books")

	books.POST("/:user_id/books", handlerBook.AddBook)
	books.PUT("/:user_id/books/:book_id", handlerBook.UpdateBook)
	books.GET("/:user_id/books/:book_id", handlerBook.GetBook)
	books.GET("/:user_id/books", handlerBook.GetAllBooks)
	books.DELETE("/:user_id/books/:book_id", handlerBook.DeleteBook)

	router.Run(port)
}
