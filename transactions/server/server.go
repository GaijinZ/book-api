package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"library/pkg/middleware"
	"library/pkg/tracing"
	"library/transactions/handler"
)

func NewRouter(transactionsHandler handler.TransactionerHandler) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1/transactions")

	v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.POST("/buy-book/:book_id",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		middleware.GetBookParam,
		transactionsHandler.BuyBook,
	)
	v1.POST("/history",
		tracing.TraceMiddleware,
		middleware.IsAuthorized,
		middleware.GetToken,
		transactionsHandler.TransactionHistory,
	)

	return router
}
