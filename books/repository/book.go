package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"library/books/models"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/pkg/utils"
)

type BookerRepository interface {
	AddBook(book *models.BookRequest) (int, error)
	UpdateBook(book *models.BookRequest) (*models.BookResponse, error)
	GetBook(id int) (*models.BookResponse, error)
	GetAllBooks() ([]models.BookResponse, error)
	DeleteBook(bookID, userID int) (int, error)
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

func (b *BookRepository) GetOrCreateAuthor(authorName string) (models.AuthorResponse, error) {
	var author models.AuthorResponse

	log := utils.GetLogger(b.ctx)

	err := b.DB.DB.QueryRow(SelectAuthor, authorName).Scan(&author.ID)
	if err != nil {
		row, err := b.DB.DB.Exec(InsertAuthor, authorName)
		if err != nil {
			log.Errorf("Failed to perform an insert query on author table: %v", err)
			return author, err
		}

		lastInsertedID, _ := row.LastInsertId()
		author.ID = int(lastInsertedID)
	}

	author.Name = authorName

	return author, nil
}

func (b *BookRepository) AddBook(book *models.BookRequest) (int, error) {
	var bookResponse models.BookResponse

	log := utils.GetLogger(b.ctx)

	isExisting, err := validateISBNExists(log, book.ISBN, b.DB.GetDB())
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
		log.Errorf("Failed to get or create author: %v", err)
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

	row, err := b.DB.DB.Exec(
		InsertBook,
		bookResponse.Name,
		bookResponse.DatePublished,
		bookResponse.ISBN,
		bookResponse.PageCount,
		bookResponse.UserID.ID,
		author.ID,
	)
	if err != nil {
		log.Errorf("Failed to perform an insert query on user book table: %d", err)
		return 0, err
	}

	lastInsertedID, _ := row.LastInsertId()

	return int(lastInsertedID), nil
}

func (b *BookRepository) UpdateBook(book *models.BookRequest) (*models.BookResponse, error) {
	var bookResponse *models.BookResponse

	log := utils.GetLogger(b.ctx)

	exists, err := postgres.CheckIDExists("books", book.ID, b.DB.GetDB())
	if err != nil {
		log.Errorf("Checking book ID error: %v", book.ID)
		return bookResponse, err
	}

	if !exists {
		log.Errorf("Book ID doesn't exists: %v", book.ID)
		errorMessage := fmt.Sprintf("Book ID doesn't exists: %v", book.ID)
		return bookResponse, errors.New(errorMessage)
	}

	isBookAssigned, err := isAssigned(log, book.ID, book.UserID.ID, b.DB.GetDB())
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
		log.Errorf("Failed to get or create author: %v", err)
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

	result, err := b.DB.DB.Exec(
		UpdateBook,
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
	bookResponse := &models.BookResponse{}

	log := utils.GetLogger(b.ctx)

	err := b.DB.DB.QueryRow(GetBook, id).Scan(
		&bookResponse.Name,
		&bookResponse.DatePublished,
		&bookResponse.ISBN,
		&bookResponse.PageCount,
		&bookResponse.Author.Name,
	)
	if err != nil {
		log.Errorf("Query row on user book table failed: %v", err)
		return bookResponse, err
	}

	return bookResponse, nil
}

func (b *BookRepository) GetAllBooks() ([]models.BookResponse, error) {
	var books []models.BookResponse

	log := utils.GetLogger(b.ctx)

	rows, err := b.DB.DB.Query(GetAllBooks)
	if err != nil {
		log.Errorf("Failed to perform a query on user book table: %d", err)
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		var book models.BookResponse

		err := rows.Scan(&book.Name, &book.DatePublished, &book.ISBN, &book.PageCount, &book.Author.Name)
		if err != nil {
			log.Errorf("Failed to scan rows: %d", err)
			return books, err
		}
		
		books = append(books, book)
	}

	return books, nil
}

func (b *BookRepository) DeleteBook(bookID, userID int) (int, error) {
	log := utils.GetLogger(b.ctx)

	exists, err := postgres.CheckIDExists("user_book", bookID, b.DB.GetDB())
	if err != nil {
		log.Errorf("Checking book ID error: %v", bookID)
		return 0, err
	}

	if !exists {
		log.Errorf("Book ID doesn't exists: %v", bookID)
		errorMessage := fmt.Sprintf("Book ID doesn't exists: %v", bookID)
		return 0, errors.New(errorMessage)
	}

	isBookAssigned, err := isAssigned(log, bookID, userID, b.DB.GetDB())
	if err != nil {
		log.Errorf("Error checking book assignment: %v", err)
		return 0, err
	}

	if !isBookAssigned {
		log.Errorf("User is not the owner of the book")
		errorMessage := "user is not the owner of the book"
		return 0, errors.New(errorMessage)
	}

	_, err = b.DB.DB.Exec(DeleteBook, bookID)
	if err != nil {
		log.Errorf("Failed to perform delete on a user book table: %d", err)
		return 0, err
	}

	return bookID, nil
}

func isAssigned(log logger.Logger, bookID, userID int, db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRow(IsAssigned, bookID, userID).Scan(&exists)
	if err != nil {
		log.Errorf("error checking book assignment: %w", err)
		return false, err
	}

	return exists, nil
}

func validateISBNExists(log logger.Logger, isbn string, db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRow(CheckISBN, isbn).Scan(&exists)
	if err != nil {
		log.Errorf("error validating ISBN %s: %w", isbn, err)
		return false, err
	}

	return exists, nil
}
