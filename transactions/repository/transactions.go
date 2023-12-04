package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"library/transactions/models"
	"time"
)

type TransactionRepository struct {
	DBPool *pgxpool.Pool
}

func NewTransactionRepository(dbPool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		DBPool: dbPool,
	}
}

func (t *TransactionRepository) BuyBook(userID, bookID, quantity int) error {
	userTransaction := &models.UserTransactionRequest{
		BookList:        &models.Book{},
		UserID:          userID,
		BookID:          bookID,
		Quantity:        quantity,
		TransactionDate: time.Now(),
	}

	tx, err := t.DBPool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	availability, err := isAvailable(userTransaction.BookID, t.DBPool)
	if err != nil {
		return err
	}

	if !availability {
		return fmt.Errorf("book is not availabe")
	}

	if err = availableQuantity(userTransaction, t.DBPool); err != nil {
		return fmt.Errorf("error checking available quantity of book")
	}

	userTransaction.BookList, err = getBookDetails(userTransaction, t.DBPool)
	if err != nil {
		return fmt.Errorf("error getting book details")
	}

	changed, err := getTransactionData(userTransaction, t.DBPool)
	if err != nil {
		tx.Rollback(context.Background())
		return fmt.Errorf("error getting transactions data")
	}

	if changed {
		return nil
	}

	newBookList, err := json.Marshal(userTransaction.BookList)
	if err != nil {
		return fmt.Errorf("failed to marshal updated book list: %v", err)
	}

	if err = updateUserTransactions(userTransaction, newBookList, t.DBPool); err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (t *TransactionRepository) TransactionHistory(userID int) ([]models.UserTransaction, error) {
	var transactions []models.UserTransaction

	query := "SELECT book_list, quantity FROM transactions WHERE user_id = $1"
	rows, err := t.DBPool.Query(context.Background(), query, userID)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		return nil, fmt.Errorf(errorMessage)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.UserTransaction

		if err = rows.Scan(&transaction.BookList, &transaction.Quantity); err != nil {
			errorMessage := "QueryRow failed: " + err.Error()
			return transactions, fmt.Errorf(errorMessage)
		}

		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		errorMessage := "Query error: " + err.Error()
		return transactions, fmt.Errorf(errorMessage)
	}

	return transactions, nil
}

// isAvailable check book for availability
func isAvailable(id int, db *pgxpool.Pool) (bool, error) {
	var availability bool

	query := "SELECT EXISTS (SELECT 1 FROM book WHERE id = $1 AND quantity > 0)"
	err := db.QueryRow(context.Background(), query, id).Scan(&availability)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		return availability, fmt.Errorf(errorMessage)
	}

	return availability, nil
}

// availableQuantity check if book has enough quantity
func availableQuantity(userTransaction *models.UserTransactionRequest, db *pgxpool.Pool) error {
	var bookQuantity int

	err := db.QueryRow(context.Background(), "SELECT quantity FROM book WHERE id = $1", userTransaction.BookID).
		Scan(&bookQuantity)
	if err != nil {
		return fmt.Errorf("failed to fetch available quantity: %v", err)
	}

	if bookQuantity < userTransaction.Quantity {
		return fmt.Errorf("not enough copies of the book available")
	}

	return nil
}

// processTransaction function to process the transaction, including the ID
func processTransaction(transaction *models.UserTransaction, userTransaction *models.UserTransactionRequest, db *pgxpool.Pool) (bool, error) {
	var book models.Book
	if err := json.Unmarshal(transaction.BookList, &book); err != nil {
		fmt.Printf("Failed to unmarshal book list for transaction %d: %v\n", transaction.ID, err)
	}

	if book.ISBN == userTransaction.BookList.ISBN {
		transactionsQuery := "UPDATE transactions SET quantity = quantity + $1 WHERE id = $2"
		_, err := db.Exec(context.Background(), transactionsQuery, userTransaction.Quantity, transaction.ID)
		if err != nil {
			return false, fmt.Errorf("failed to update available quantity in the table: %v", err)
		}

		bookQuery := "UPDATE book SET quantity = quantity - $1 WHERE id = $2"
		_, err = db.Exec(context.Background(), bookQuery, userTransaction.Quantity, userTransaction.BookList.ID)
		if err != nil {
			return false, fmt.Errorf("failed to update available quantity in the table: %v", err)
		}

	}

	return true, nil
}

func getBookDetails(userTransaction *models.UserTransactionRequest, db *pgxpool.Pool) (*models.Book, error) {
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
		return userTransaction.BookList, fmt.Errorf("failed to fetch book details: %v", err)
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
			return userTransaction.BookList, fmt.Errorf("failed to scan book details: %v", err)
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

func getTransactionData(userTransaction *models.UserTransactionRequest, db *pgxpool.Pool) (bool, error) {
	var changed bool

	query := "SELECT id, book_list FROM transactions WHERE user_id = $1"
	rows, err := db.Query(context.Background(), query, userTransaction.UserID)
	if err != nil {
		return changed, fmt.Errorf("failed to fetch existing transaction data: %v", err)
	}

	for rows.Next() {
		var transaction models.UserTransaction

		if err = rows.Scan(&transaction.ID, &transaction.BookList); err != nil {
			return changed, fmt.Errorf("failed to scan row: %v", err)
		}

		// check and update transaction quantity
		changed, err = processTransaction(&transaction, userTransaction, db)
		if err != nil {
			return changed, nil
		}

	}

	if err = rows.Err(); err != nil {
		return changed, fmt.Errorf("error iterating through rows: %v", err)
	}

	return changed, nil
}

func updateUserTransactions(userTransaction *models.UserTransactionRequest, newBookList []byte, db *pgxpool.Pool) error {
	transactionsQuery := "INSERT INTO transactions (user_id, book_list, quantity, transaction_date) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(
		context.Background(),
		transactionsQuery,
		userTransaction.UserID,
		newBookList, userTransaction.Quantity,
		userTransaction.TransactionDate,
	)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %v", err)
	}

	bookQuery := "UPDATE book SET quantity = quantity - $1 WHERE id = $2"
	_, err = db.Exec(context.Background(), bookQuery, userTransaction.Quantity, userTransaction.BookList.ID)
	if err != nil {
		return fmt.Errorf("failed to update available quantity in the table: %v", err)
	}

	return nil
}
