package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/pkg/middleware"
	"library/pkg/postgres"
	"library/pkg/tracing"
	"library/transactions/handler"
	"library/transactions/repository"
)

func Run(ctx *context.Context, port string) {
	dbPool := postgres.GetConnection()
	defer dbPool.Close()

	transactionsRepository := repository.NewTransactionRepository(*ctx, dbPool)
	transactionsHandler := handler.NewTransactionHandler(*ctx, transactionsRepository)

	router := gin.Default()

	v1 := router.Group("/v1/transactions")

	v1.POST("/buy-book/:book_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		transactionsHandler.BuyBook,
	)
	v1.POST("/transactions",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		transactionsHandler.TransactionHistory,
	)

	router.Run(port)
}
