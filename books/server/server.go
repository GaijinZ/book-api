package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"library/books/handler"
	"library/books/repository"
	"library/pkg/config"
	"library/pkg/middleware"
	"library/pkg/postgres"
	"library/pkg/tracing"
	"library/pkg/utils"

	"time"
)

func Run(ctx *context.Context, cfg config.GlobalEnv) {
	log := utils.GetLogger(*ctx)

	configDB := postgres.DBConfig{
		DriverName:      "postgres",
		DataSourceName:  cfg.PostgresBooks,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}

	db, err := postgres.NewDB(*ctx, configDB)
	if err != nil {
		log.Errorf("Failed to configure db connection: %v", err)
	}
	defer db.Close()

	bookRepository := repository.NewBookRepository(*ctx, *db)
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

	router.Run(":" + cfg.BooksServerPort)
}
