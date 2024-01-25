package repository

import (
	"context"
	"encoding/json"

	"library/pkg/logger"
	"library/pkg/postgres"
	"library/pkg/utils"
	"library/transactions/models"

	"time"
)

type TransactionerRepository interface {
	BuyBook(userID, bookID, quantity int) (int, error)
	TransactionHistory(userID int) ([]models.UserTransactionResponse, error)
}

type TransactionRepository struct {
	ctx context.Context
	DB  postgres.DB
}

func NewTransactionRepository(ctx context.Context, db postgres.DB) TransactionerRepository {
	return &TransactionRepository{
		ctx: ctx,
		DB:  db,
	}
}

func (t *TransactionRepository) BuyBook(userID, bookID, quantity int) (int, error) {
	log := utils.GetLogger(t.ctx)

	userTransaction := &models.UserTransactionRequest{
		BookList:        &models.Book{},
		UserID:          userID,
		BookID:          bookID,
		Quantity:        quantity,
		TransactionDate: time.Now(),
	}

	tx, err := t.DB.DB.Begin()
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return 0, err
	}

	availability, err := isAvailable(log, userTransaction.BookID, t.DB)
	if err != nil {
		log.Errorf("Failed to check book availbility: %v", err)
		return 0, err
	}

	if !availability {
		log.Errorf("Book is not available, book id: %v", bookID)
		return 0, err
	}

	if err = availableQuantity(log, userTransaction, t.DB); err != nil {
		log.Errorf("Failed to check available quantity of book")
		return 0, err
	}

	userTransaction.BookList, err = getBookDetails(log, userTransaction, t.DB)
	if err != nil {
		log.Errorf("Failed to get book details: %v", err)
		return 0, err
	}

	changed, err := getTransactionData(log, userTransaction, t.DB)
	if err != nil {
		log.Errorf("Failed to get transaction data")
		tx.Rollback()
		return 0, err
	}

	if changed == true {
		return 0, nil
	}

	newBookList, err := json.Marshal(userTransaction.BookList)
	if err != nil {
		log.Errorf("Failed to marshal updated book list: %v", err)
		return 0, err
	}

	if err = updateUserTransactions(log, userTransaction, newBookList, t.DB); err != nil {
		log.Errorf("Failed to updated user transactions: %v", err)
		return 0, err
	}
	log.Infof("user transactions updated")

	err = tx.Commit()
	if err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return 0, err
	}

	return bookID, nil
}

func (t *TransactionRepository) TransactionHistory(userID int) ([]models.UserTransactionResponse, error) {
	var transactions []models.UserTransactionResponse
	log := utils.GetLogger(t.ctx)

	query := "SELECT book_list, quantity FROM transactions WHERE user_id = $1"
	rows, err := t.DB.DB.Query(query, userID)
	if err != nil {
		log.Errorf("Failed to query row: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.UserTransactionResponse

		if err = rows.Scan(&transaction.BookList, &transaction.Quantity); err != nil {
			log.Errorf("Failed to scan into UserTransaction: %v", err)
			return transactions, err
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		log.Errorf("Query failed: %v", err)
		return transactions, err
	}

	return transactions, nil
}

// isAvailable check book for availability
func isAvailable(log logger.Logger, id int, db postgres.DB) (bool, error) {
	var availability bool

	query := "SELECT EXISTS (SELECT 1 FROM book WHERE id = $1 AND quantity > 0)"
	err := db.DB.QueryRow(query, id).Scan(&availability)
	if err != nil {
		log.Errorf("Failed to query row: %v", err)
		return availability, err
	}

	return availability, nil
}

// availableQuantity check if book has enough quantity
func availableQuantity(log logger.Logger, userTransaction *models.UserTransactionRequest, db postgres.DB) error {
	var bookQuantity int

	err := db.DB.QueryRow("SELECT quantity FROM book WHERE id = $1", userTransaction.BookID).
		Scan(&bookQuantity)
	if err != nil {
		log.Errorf("Failed to fetch available quantity: %v", err)
		return err
	}

	if bookQuantity < userTransaction.Quantity {
		log.Errorf("Not enough copies available")
		return err
	}

	return nil
}

// processTransaction function to process the transaction, including the ID
func processTransaction(
	log logger.Logger,
	transaction *models.UserTransactionResponse,
	userTransaction *models.UserTransactionRequest,
	db postgres.DB,
) (bool, error) {
	var book models.Book
	if err := json.Unmarshal(transaction.BookList, &book); err != nil {
		log.Errorf("Failed to unmarshal book list for transaction %d: %v", transaction.ID, err)
		return false, err
	}

	if book.ISBN == userTransaction.BookList.ISBN {
		transactionsQuery := "UPDATE transactions SET quantity = quantity + $1 WHERE id = $2"
		_, err := db.DB.Exec(transactionsQuery, userTransaction.Quantity, transaction.ID)
		if err != nil {
			log.Errorf("Failed to update available quantity: %v", err)
			return false, err
		}

		bookQuery := "UPDATE book SET quantity = quantity - $1 WHERE id = $2"
		_, err = db.DB.Exec(bookQuery, userTransaction.Quantity, userTransaction.BookList.ID)
		if err != nil {
			log.Errorf("Failed to update available quantity: %v", err)
			return false, err
		}

	}

	return true, nil
}

func getBookDetails(
	log logger.Logger,
	userTransaction *models.UserTransactionRequest,
	db postgres.DB,
) (*models.Book, error) {
	var mapID int

	queryBook := `
	   SELECT
	       book.id,
	       book.name,
	       book.date_published,
	       book.isbn,
	       book.page_count,
	       author.name AS author_name
	   FROM
	       book
	       JOIN book_authors ON book.id = book_authors.book_id
	       JOIN author ON author.id = book_authors.author_id
	   WHERE
	       book.id = $1;
	`

	rows, err := db.DB.Query(queryBook, userTransaction.BookID)
	if err != nil {
		log.Errorf("Failed to fetch book details: %v", err)
		return userTransaction.BookList, err
	}
	defer rows.Close()

	//// map to store book details indexed by book ID
	bookDetails := make(map[int]*models.Book)

	for rows.Next() {
		var authorName string

		if err = rows.Scan(
			&userTransaction.BookList.ID,
			&userTransaction.BookList.Name,
			&userTransaction.BookList.DatePublished,
			&userTransaction.BookList.ISBN,
			&userTransaction.BookList.PageCount,
			&authorName,
		); err != nil {
			log.Errorf("Failed to fetch book details: %v", err)
			return userTransaction.BookList, err
		}

		// Check if the book ID is already in the map
		existingBook, ok := bookDetails[userTransaction.BookList.ID]
		if !ok {
			// If not, create a new Book with the current details
			existingBook = &models.Book{
				ID:            userTransaction.BookList.ID,
				Name:          userTransaction.BookList.Name,
				DatePublished: userTransaction.BookList.DatePublished,
				ISBN:          userTransaction.BookList.ISBN,
				PageCount:     userTransaction.BookList.PageCount,
				Author:        []models.Author{},
			}
		}

		// Append the new author to the existing or new Book
		existingBook.Author = append(existingBook.Author, models.Author{Name: authorName})
		bookDetails[userTransaction.BookList.ID] = existingBook
		mapID = userTransaction.BookList.ID
	}

	return bookDetails[mapID], nil
}

func getTransactionData(
	log logger.Logger,
	userTransaction *models.UserTransactionRequest,
	db postgres.DB,
) (bool, error) {
	var changed bool

	query := "SELECT id, book_list FROM transactions WHERE user_id = $1"
	rows, err := db.DB.Query(query, userTransaction.UserID)
	if err != nil {
		log.Errorf("Failed to fetch existing transactions: %v", err)
		return changed, err
	}

	for rows.Next() {
		var transaction models.UserTransactionResponse

		if err = rows.Scan(&transaction.ID, &transaction.BookList); err != nil {
			log.Errorf("Failed to scan rows: %v", err)
			return changed, err
		}

		changed, err = processTransaction(log, &transaction, userTransaction, db)
		if err != nil {
			return changed, nil
		}

	}

	if err = rows.Err(); err != nil {
		log.Errorf("Failed to iterate through rows: %v", err)
		return changed, err
	}

	return changed, nil
}

func updateUserTransactions(
	log logger.Logger,
	userTransaction *models.UserTransactionRequest,
	newBookList []byte,
	db postgres.DB,
) error {
	transactionsQuery := "INSERT INTO transactions (user_id, book_list, quantity, transaction_date) VALUES ($1, $2, $3, $4)"
	_, err := db.DB.Exec(
		transactionsQuery,
		userTransaction.UserID,
		newBookList,
		userTransaction.Quantity,
		userTransaction.TransactionDate,
	)
	if err != nil {
		log.Errorf("Failed to update transactions: %v", err)
		return err
	}

	bookQuery := "UPDATE book SET quantity = quantity - $1 WHERE id = $2"
	_, err = db.DB.Exec(bookQuery, userTransaction.Quantity, userTransaction.BookList.ID)
	if err != nil {
		log.Errorf("Failed to update available quantity: %v", err)
		return err
	}

	return nil
}
