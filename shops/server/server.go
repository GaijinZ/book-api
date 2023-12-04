package server

import (
	"github.com/gin-gonic/gin"
	"library/pkg/postgres"
	"library/shops/handler"
	"library/shops/repository"
)

func Run(port string) {
	dbPool := postgres.GetConnection()
	defer dbPool.Close()

	shopRepository := repository.NewShopRepository(dbPool)
	shopHandler := handler.NewShopHandler(shopRepository)

	router := gin.Default()

	v1 := router.Group("/v1/shops")

	v1.POST("/load-books", shopHandler.LoadBooks)

	router.Run(port)
}
