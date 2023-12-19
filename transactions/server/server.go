package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/pkg/config"
	"library/pkg/middleware"
	"library/pkg/postgres"
	"library/pkg/tracing"
	"library/pkg/utils"
	"library/transactions/handler"
	"library/transactions/repository"
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

	transactionsRepository := repository.NewTransactionRepository(*ctx, *db)
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

	router.Run(":" + cfg.TransactionsServerPort)
}
