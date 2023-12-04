package server

import (
	"github.com/gin-gonic/gin"
	"library/pkg/postgres"
	"library/transactions/handler"
	"library/transactions/repository"
)

func Run(port string) {
	dbPool := postgres.GetConnection()
	defer dbPool.Close()

	transactionsRepository := repository.NewTransactionRepository(dbPool)
	transactionsHandler := handler.NewTransactionHandler(transactionsRepository)

	router := gin.Default()

	v1 := router.Group("/v1/transactions")

	v1.POST("/buy-book/:book_id", transactionsHandler.BuyBook)
	v1.POST("/transactions", transactionsHandler.TransactionHistory)

	router.Run(port)
}
