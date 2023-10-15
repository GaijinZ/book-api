package server

import (
	"github.com/gin-gonic/gin"
	"library/shops/handler"
	"library/shops/postgres"
	"library/shops/repository"
)

func Run(port string) {
	dbPool := postgres.GetPostgresConnectionString()
	defer dbPool.Close()

	shopRepository := repository.NewShopRepository(dbPool)
	shopHandler := handler.NewShopHandler(shopRepository)

	router := gin.Default()

	v1 := router.Group("/v1")

	v1.POST("/shop-list", shopHandler.LoadBooks)

	router.Run(port)
}
