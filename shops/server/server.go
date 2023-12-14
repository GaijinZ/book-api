package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/pkg/postgres"
	"library/shops/handler"
	"library/shops/repository"
)

func Run(ctx *context.Context, port string) {
	dbPool := postgres.GetConnection()
	defer dbPool.Close()

	shopRepository := repository.NewShopRepository(*ctx, dbPool)
	shopHandler := handler.NewShopHandler(*ctx, shopRepository)

	router := gin.Default()

	v1 := router.Group("/v1/shops")

	v1.POST("/load-books", shopHandler.LoadBooks)

	router.Run(port)
}
