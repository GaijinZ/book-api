package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/pkg/logger"
	"library/shops/repository"
	"net/http"
)

type ShopperHandler interface {
	LoadBooks(c *gin.Context)
}

type ShopHandler struct {
	ctx            context.Context
	shopRepository repository.ShopperRepository
}

func NewShopHandler(ctx context.Context, shopperRepository repository.ShopperRepository) ShopperHandler {
	return &ShopHandler{
		ctx:            ctx,
		shopRepository: shopperRepository,
	}
}

func (s *ShopHandler) LoadBooks(c *gin.Context) {
	log := s.ctx.Value("logger").(logger.Logger)
	err := s.shopRepository.LoadBooks()
	if err != nil {
		errorMessage := "Count not fetch books: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	log.Infof("Available shops:")
	log.Infof("Hobbit")
	log.Infof("War")
	log.Infof("Lego")
}
