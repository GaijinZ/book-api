package repository

import (
	"context"
	"errors"
	"fmt"
	"library/books/models"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/pkg/utils"
	"time"
)

type BookerRepository interface {
	AddBook(book *models.BookRequest) (int, error)
	UpdateBook(book *models.BookRequest) (*models.BookResponse, error)
	GetBook(id int) (*models.BookResponse, error)
	GetAllBooks() ([]models.BookResponse, error)
	DeleteBook(bookID, userID int) (int, error)
	GetDBPool() postgres.DB
}

type BookRepository struct {
	ctx context.Context
	DB  postgres.DB
}

func NewBookRepository(ctx context.Context, db postgres.DB) BookerRepository {
	return &BookRepository{
		ctx: ctx,
		DB:  db,
	}
}

func (b *BookRepository) GetDBPool() postgres.DB {
	return b.DB
}

func (b *BookRepository) GetOrCreateAuthor(authorName string) (models.AuthorResponse, error) {
	var author models.AuthorResponse

	log := utils.GetLogger(b.ctx)

	query := "SELECT id FROM author WHERE name = $1"
	err := b.DB.DB.QueryRow(query, authorName).Scan(&author.ID)
	if err != nil {
		query = "INSERT INTO author (name) VALUES ($1) RETURNING id"
		err = b.DB.DB.QueryRow(query, authorName).Scan(&author.ID)
		if err != nil {
			log.Errorf("Failed to perform an insert query on author table: %d", err)
			return author, err
		}
	}

	author.Name = authorName

	return author, nil
}

func (b *BookRepository) AddBook(book *models.BookRequest) (int, error) {
	var bookResponse models.BookResponse

	log := utils.GetLogger(b.ctx)

	isExisting, err := validateISBNExists(log, book.ISBN, b.GetDBPool())
	if err != nil {
		log.Errorf("ISBN validation error: %v", err)
		return 0, err
	}

	if isExisting {
		log.Errorf("ISBN already exists")
		errorMessage := "book with this ISBN already exists"
		return 0, errors.New(errorMessage)
	}

	author, err := b.GetOrCreateAuthor(book.Author.Name)
	if err != nil {
		log.Errorf("Failed to get or create author: %d", err)
		return 0, err
	}

	bookResponse = models.BookResponse{
		Name:          book.Name,
		DatePublished: book.DatePublished,
		ISBN:          book.ISBN,
		PageCount:     book.PageCount,
		UserID:        book.UserID,
		Author:        author,
	}

	bookResponse.Author.ID = author.ID

	var lastInsertedID int
	insertQuery := `
INSERT INTO user_book (name, date_published, isbn, page_count, user_id, author_id) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING id
`
	err = b.DB.DB.QueryRow(
		insertQuery,
		bookResponse.Name,
		bookResponse.DatePublished,
		bookResponse.ISBN,
		bookResponse.PageCount,
		bookResponse.UserID.ID,
		author.ID,
	).Scan(&lastInsertedID)
	if err != nil {
		log.Errorf("Failed to perform an insert query on user book table: %d", err)
		return 0, err
	}

	return lastInsertedID, nil
}

func (b *BookRepository) UpdateBook(book *models.BookRequest) (*models.BookResponse, error) {
	var bookResponse *models.BookResponse

	log := utils.GetLogger(b.ctx)

	exists, err := postgres.CheckIDExists("books", book.ID, b.GetDBPool())
	if err != nil {
		log.Errorf("Checking book ID error: %v", book.ID)
		return bookResponse, err
	}

	if !exists {
		log.Errorf("Book ID doesn't exists: %v", book.ID)
		errorMessage := fmt.Sprintf("Book ID doesn't exists: %v", book.ID)
		return bookResponse, errors.New(errorMessage)
	}

	isBookAssigned, err := isAssigned(log, book.ID, book.UserID.ID, b.GetDBPool())
	if err != nil {
		log.Errorf("Error checking book assignment: %v", err)
		return bookResponse, err
	}

	if !isBookAssigned {
		log.Errorf("User is not the owner of the book")
		errorMessage := "user is not the owner of the book"
		return bookResponse, errors.New(errorMessage)
	}

	author, err := b.GetOrCreateAuthor(book.Author.Name)
	if err != nil {
		log.Errorf("Failed to get or create author: %d", err)
		return bookResponse, err
	}

	bookResponse = &models.BookResponse{
		ID:            book.ID,
		Name:          book.Name,
		DatePublished: book.DatePublished,
		ISBN:          book.ISBN,
		PageCount:     book.PageCount,
		Author: models.AuthorResponse{
			ID:   author.ID,
			Name: book.Author.Name,
		},
	}

	updateQuery := "UPDATE user_book SET name=$1, date_published=$2, isbn=$3, page_count=$4, author_id=$5 WHERE id=$6"
	result, err := b.DB.DB.Exec(
		updateQuery,
		bookResponse.Name,
		bookResponse.DatePublished,
		bookResponse.ISBN,
		bookResponse.PageCount,
		bookResponse.Author.ID,
		bookResponse.ID,
	)
	if err != nil {
		log.Errorf("Failed to perform an update query in user book table: %d", err)
		return bookResponse, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Errorf("Error number of rows affected: %v", err)
		return bookResponse, err
	}

	if affectedRows == 0 {
		log.Errorf("No rows were affected: %d", err)
		return bookResponse, err
	}

	return bookResponse, nil
}

func (b *BookRepository) GetBook(id int) (*models.BookResponse, error) {
	var date time.Time
	var bookResponse *models.BookResponse

	log := utils.GetLogger(b.ctx)

	getQuery := `
	SELECT b.name, b.date_published, b.isbn, b.page_count, a.name
	FROM user_book AS b
	JOIN author AS a ON b.author_id = a.id
	WHERE b.id = $1
`
	err := b.DB.DB.QueryRow(getQuery, id).Scan(
		&bookResponse.Name, &date, &bookResponse.ISBN, &bookResponse.PageCount, &bookResponse.Author.Name,
	)
	if err != nil {
		log.Errorf("Query row on user book table failed: %d", err)
		return bookResponse, err
	}

	bookResponse.DatePublished = date.Format("2006-01-02")

	return bookResponse, nil
}

func (b *BookRepository) GetAllBooks() ([]models.BookResponse, error) {
	var books []models.BookResponse
	var date time.Time

	log := utils.GetLogger(b.ctx)

	getAllQuery := `
	SELECT b.name, b.date_published, b.isbn, b.page_count, a.name
	FROM user_book AS b
	JOIN author AS a ON b.author_id = a.id
	ORDER BY a.id
	`
	rows, err := b.DB.DB.Query(getAllQuery)
	if err != nil {
		log.Errorf("Failed to perform a query on user book table: %d", err)
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		var book models.BookResponse

		err := rows.Scan(&book.Name, &date, &book.ISBN, &book.PageCount, &book.Author.Name)
		if err != nil {
			log.Errorf("Failed to scan rows: %d", err)
			return books, err
		}

		book.DatePublished = date.Format("2006-01-02")
		books = append(books, book)
	}

	return books, nil
}

func (b *BookRepository) DeleteBook(bookID, userID int) (int, error) {
	log := utils.GetLogger(b.ctx)

	exists, err := postgres.CheckIDExists("user_book", bookID, b.GetDBPool())
	if err != nil {
		log.Errorf("Checking book ID error: %v", bookID)
		return 0, err
	}

	if !exists {
		log.Errorf("Book ID doesn't exists: %v", bookID)
		errorMessage := fmt.Sprintf("Book ID doesn't exists: %v", bookID)
		return 0, errors.New(errorMessage)
	}

	isBookAssigned, err := isAssigned(log, bookID, userID, b.GetDBPool())
	if err != nil {
		log.Errorf("Error checking book assignment: %v", err)
		return 0, err
	}

	if !isBookAssigned {
		log.Errorf("User is not the owner of the book")
		errorMessage := "user is not the owner of the book"
		return 0, errors.New(errorMessage)
	}

	query := "DELETE FROM user_book WHERE id=$1"
	_, err = b.DB.DB.Exec(query, bookID)
	if err != nil {
		log.Errorf("Failed to perform delete on a user book table: %d", err)
		return 0, err
	}

	return bookID, nil
}

func isAssigned(log logger.Logger, bookID, userID int, db postgres.DB) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM user_book WHERE id = $1 AND user_id = $2)"

	var exists bool
	err := db.DB.QueryRow(query, bookID, userID).Scan(&exists)
	if err != nil {
		log.Errorf("error checking book assignment: %w", err)
		return false, err
	}

	return exists, nil
}

func validateISBNExists(log logger.Logger, isbn string, db postgres.DB) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM user_book WHERE isbn = $1)"

	var exists bool
	err := db.DB.QueryRow(query, isbn).Scan(&exists)
	if err != nil {
		log.Errorf("error validating ISBN %s: %w", isbn, err)
		return false, err
	}

	return exists, nil
}
