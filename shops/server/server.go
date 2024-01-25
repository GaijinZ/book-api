package server

import (
	"github.com/gin-gonic/gin"
	"library/shops/handler"
)

func NewRouter(shopHandler handler.ShopperHandler) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1/shops")

	v1.POST("/load-books", shopHandler.LoadBooks)

	return router
}
