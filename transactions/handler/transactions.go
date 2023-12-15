package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/transactions/models"
	"library/transactions/repository"
	"net/http"
	"strconv"
)

type TransactionerHandler interface {
	BuyBook(*gin.Context)
	TransactionHistory(*gin.Context)
}

type TransactionHandler struct {
	ctx                   context.Context
	transactionRepository repository.TransactionerRepository
}

func NewTransactionHandler(ctx context.Context, transactionHandler repository.TransactionerRepository) TransactionerHandler {
	return &TransactionHandler{
		ctx:                   ctx,
		transactionRepository: transactionHandler,
	}
}

func (t *TransactionHandler) BuyBook(c *gin.Context) {
	transaction := models.TransactionResponse{}

	userID := c.GetInt("userID")

	bookID, err := strconv.Atoi(c.Param("book_id")) // move getting param to midleware, and pass it as context
	if err != nil {
		errorMessage := "Wrong book ID: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = t.transactionRepository.BuyBook(userID, bookID, transaction.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Book added successfully": "book"})
}

func (t *TransactionHandler) TransactionHistory(c *gin.Context) {
	userID := c.GetInt("userID")

	transactions, err := t.transactionRepository.TransactionHistory(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
