package server

import (
	"library/database/postgres"
	"library/internal/books/handler"
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

	v1.POST("/user/:user_id/book", handlerBook.AddBook)
	v1.PUT("/user/:user_id/book/:book_id", handlerBook.UpdateBook)
	v1.GET("/user/:user_id/book/:book_id", handlerBook.GetBook)
	v1.GET("/user/:user_id/book", handlerBook.GetAllBooks)
	v1.DELETE("/user/:user_id/book/:book_id", handlerBook.DeleteBook)

	router.Run(port)
}
