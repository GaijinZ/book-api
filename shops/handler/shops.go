package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library/shops/downloads_api"
	"library/shops/repository"
)

type ShopHandler struct {
	shopRepository *repository.ShopRepository
}

func NewShopHandler(shopRepository *repository.ShopRepository) *ShopHandler {
	return &ShopHandler{
		shopRepository: shopRepository,
	}
}

func (s *ShopHandler) LoadBooks(c *gin.Context) {
	booksResponse, err := downloads_api.GetBooks("Hobbit")
	if err != nil {
		fmt.Println("Error fetching books:", err)
		return
	}

	err = s.shopRepository.LoadBooks(booksResponse)
	if err != nil {
		return
	}

	fmt.Println("Available shops:")
	fmt.Println("Hobbit")
}
