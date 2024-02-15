package handler

import (
	"context"

	"library/pkg/utils"
	"library/transactions/models"
	"library/transactions/repository"
	"net/http"

	"github.com/gin-gonic/gin"
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

// BuyBook performs a transaction to buy a book.
//
//	@Summary		Buy a book
//	@Description	Buy a book with the provided transaction data
//	@Accept			json
//	@Produce		json
//	@Param			bookID		path		int											true	"Book ID"
//	@Param			transaction	body		models.TransactionResponse							true	"Transaction data"
//	@Success		201
//	@Failure		400
//	@Failure		500
//	@Router			/v1/transactions/buy-book/{bookID} [post]
func (t *TransactionHandler) BuyBook(c *gin.Context) {
	transaction := models.TransactionResponse{}
	log := utils.GetLogger(t.ctx)

	userID := 1
	bookID := c.GetInt("bookID")

	if err := c.ShouldBindJSON(&transaction); err != nil {
		log.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := t.transactionRepository.BuyBook(userID, bookID, transaction.Quantity)
	if err != nil {
		log.Errorf("Transaction repository error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Transaction successful: %v", id)
	c.JSON(http.StatusCreated, gin.H{"Book added successfully": "book"})
}

// TransactionHistory retrieves the transaction history for a user.
//
//	@Summary		Get transaction history
//	@Description	Retrieves the transaction history for the specified user
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int										true	"User ID"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/v1/transactions/history/{userID} [get]
func (t *TransactionHandler) TransactionHistory(c *gin.Context) {
	log := utils.GetLogger(t.ctx)
	userID := c.GetInt("userID")

	transactions, err := t.transactionRepository.TransactionHistory(userID)
	if err != nil {
		log.Errorf("Transaction repository error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Transaction history acuired for user: %v", userID)
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
