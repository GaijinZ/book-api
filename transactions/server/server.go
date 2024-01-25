package server

import (
	"github.com/gin-gonic/gin"
	"library/pkg/middleware"
	"library/pkg/tracing"
	"library/transactions/handler"
)

func NewRouter(transactionsHandler handler.TransactionerHandler) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1/transactions")

	v1.POST("/buy-book/:book_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		middleware.GetBookParam,
		transactionsHandler.BuyBook,
	)
	v1.POST("/transactions",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		transactionsHandler.TransactionHistory,
	)

	return router
}
