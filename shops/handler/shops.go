package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library/shops/repository"
	"net/http"
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
	err := s.shopRepository.LoadBooks()
	if err != nil {
		errorMessage := "Count not fetch books: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	fmt.Println("Available shops:")
	fmt.Println("Hobbit")
	fmt.Println("War")
	fmt.Println("Lego")
}
