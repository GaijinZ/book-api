package handler

import (
	"github.com/gin-gonic/gin"
	"library/transactions/models"
	"library/transactions/repository"
	"net/http"
	"strconv"
)

type Transaction interface {
	BuyBook(*gin.Context)
	TransactionHistory(*gin.Context)
}

type TransactionHandler struct {
	transactionRepository *repository.TransactionRepository
}

//func NewTransactionHandler(transactionHandler *repository.TransactionRepository) *TransactionHandler {
//	return &TransactionHandler{
//		transactionRepository: transactionHandler,
//	}
//}

func NewTransactionHandler(transactionHandler *repository.TransactionRepository) Transaction {
	tr := TransactionHandler{}
	tr.transactionRepository = transactionHandler

	return &tr
}

func (t *TransactionHandler) BuyBook(c *gin.Context) {
	transaction := models.TransactionResponse{}

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		errorMessage := "Wrong book ID: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//token, err := c.Cookie("token")
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//claims, err := middleware.VerifyJWT(token)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//	return
	//}
	//
	//userID, err := strconv.Atoi(claims.UserID)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	if err = t.transactionRepository.BuyBook(1, bookID, transaction.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Book added successfully": "book"})
}

func (t *TransactionHandler) TransactionHistory(c *gin.Context) {
	//token, err := c.Cookie("token")
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	//
	//claims, err := middleware.VerifyJWT(token)
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	//	return
	//}
	//
	//userID, err := strconv.Atoi(claims.UserID)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}

	transactions, err := t.transactionRepository.TransactionHistory(1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
