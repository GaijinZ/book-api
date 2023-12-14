package repository

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"library/pkg/logger"
	"library/transactions/models"
	"time"
)

type TransactionerRepository interface {
	BuyBook(userID, bookID, quantity int) error
	TransactionHistory(userID int) ([]models.UserTransaction, error)
}

type TransactionRepository struct {
	ctx    context.Context
	DBPool *pgxpool.Pool
}

func NewTransactionRepository(ctx context.Context, dbPool *pgxpool.Pool) TransactionerRepository {
	return &TransactionRepository{
		ctx:    ctx,
		DBPool: dbPool,
	}
}

func (t *TransactionRepository) BuyBook(userID, bookID, quantity int) error {
	log := t.ctx.Value("logger").(logger.Logger)
	userTransaction := &models.UserTransactionRequest{
		BookList:        &models.Book{},
		UserID:          userID,
		BookID:          bookID,
		Quantity:        quantity,
		TransactionDate: time.Now(),
	}

	tx, err := t.DBPool.Begin(context.Background())
	if err != nil {
		log.Errorf("Failed to begin transaction: %v", err)
		return err
	}

	availability, err := isAvailable(log, userTransaction.BookID, t.DBPool)
	if err != nil {
		return err
	}

	if !availability {
		log.Errorf("Book is not available")
		return err
	}

	if err = availableQuantity(log, userTransaction, t.DBPool); err != nil {
		log.Errorf("Failed to check available quantity of book")
		return err
	}

	userTransaction.BookList, err = getBookDetails(log, userTransaction, t.DBPool)
	if err != nil {
		log.Errorf("Failed to get book details: %v", err)
		return err
	}

	changed, err := getTransactionData(log, userTransaction, t.DBPool)
	if err != nil {
		log.Errorf("Failed to get transaction data")
		tx.Rollback(context.Background())
		return err
	}

	if changed {
		return nil
	}

	newBookList, err := json.Marshal(userTransaction.BookList)
	if err != nil {
		log.Errorf("Failed to marshal updated book list: %v", err)
		return err
	}

	if err = updateUserTransactions(log, userTransaction, newBookList, t.DBPool); err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Errorf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (t *TransactionRepository) TransactionHistory(userID int) ([]models.UserTransaction, error) {
	var transactions []models.UserTransaction
	log := t.ctx.Value("logger").(logger.Logger)

	query := "SELECT book_list, quantity FROM transactions WHERE user_id = $1"
	rows, err := t.DBPool.Query(context.Background(), query, userID)
	if err != nil {
		log.Errorf("Failed to query row: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.UserTransaction

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
func isAvailable(log logger.Logger, id int, db *pgxpool.Pool) (bool, error) {
	var availability bool

	query := "SELECT EXISTS (SELECT 1 FROM book WHERE id = $1 AND quantity > 0)"
	err := db.QueryRow(context.Background(), query, id).Scan(&availability)
	if err != nil {
		log.Errorf("Failed to query row: %v", err)
		return availability, err
	}

	return availability, nil
}

// availableQuantity check if book has enough quantity
func availableQuantity(log logger.Logger, userTransaction *models.UserTransactionRequest, db *pgxpool.Pool) error {
	var bookQuantity int

	err := db.QueryRow(context.Background(), "SELECT quantity FROM book WHERE id = $1", userTransaction.BookID).
		Scan(&bookQuantity)
	if err != nil {
		log.Errorf("Failed to fetch available quantity: %V", err)
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
	transaction *models.UserTransaction,
	userTransaction *models.UserTransactionRequest,
	db *pgxpool.Pool,
) (bool, error) {
	var book models.Book
	if err := json.Unmarshal(transaction.BookList, &book); err != nil {
		log.Errorf("Failed to unmarshal book list for transaction %d: %v", transaction.ID, err)
		return false, err
	}

	if book.ISBN == userTransaction.BookList.ISBN {
		transactionsQuery := "UPDATE transactions SET quantity = quantity + $1 WHERE id = $2"
		_, err := db.Exec(context.Background(), transactionsQuery, userTransaction.Quantity, transaction.ID)
		if err != nil {
			log.Errorf("Failed to update available quantity: %v", err)
			return false, err
		}

		bookQuery := "UPDATE book SET quantity = quantity - $1 WHERE id = $2"
		_, err = db.Exec(context.Background(), bookQuery, userTransaction.Quantity, userTransaction.BookList.ID)
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
	db *pgxpool.Pool,
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

	rows, err := db.Query(context.Background(), queryBook, userTransaction.BookID)
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
	db *pgxpool.Pool,
) (bool, error) {
	var changed bool

	query := "SELECT id, book_list FROM transactions WHERE user_id = $1"
	rows, err := db.Query(context.Background(), query, userTransaction.UserID)
	if err != nil {
		log.Errorf("Failed to fetch existing transactions: %v", err)
		return changed, err
	}

	for rows.Next() {
		var transaction models.UserTransaction

		if err = rows.Scan(&transaction.ID, &transaction.BookList); err != nil {
			log.Errorf("Failed to scan rows: %v", err)
			return changed, err
		}

		// check and update transaction quantity
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
	db *pgxpool.Pool,
) error {
	transactionsQuery := "INSERT INTO transactions (user_id, book_list, quantity, transaction_date) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(
		context.Background(),
		transactionsQuery,
		userTransaction.UserID,
		newBookList, userTransaction.Quantity,
		userTransaction.TransactionDate,
	)
	if err != nil {
		log.Errorf("Failed to update transactions: %v", err)
		return err
	}

	bookQuery := "UPDATE book SET quantity = quantity - $1 WHERE id = $2"
	_, err = db.Exec(context.Background(), bookQuery, userTransaction.Quantity, userTransaction.BookList.ID)
	if err != nil {
		log.Errorf("Failed to update available quantity: %v", err)
		return err
	}

	return nil
}
