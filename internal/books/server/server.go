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

	v1.POST("/users/:user_id/books", handlerBook.AddBook)
	v1.PUT("/users/:user_id/books/:book_id", handlerBook.UpdateBook)
	v1.GET("/users/:user_id/books/:book_id", handlerBook.GetBook)
	v1.GET("/users/:user_id/books", handlerBook.GetAllBooks)
	v1.DELETE("/users/:user_id/books/:book_id", handlerBook.DeleteBook)

	router.Run(port)
}
