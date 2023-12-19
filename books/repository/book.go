package repository

import (
	"context"

	"library/books/models"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/pkg/utils"

	"time"
)

type BookerRepository interface {
	AddBook(book *models.BookRequest) error
	UpdateBook(book *models.BookRequest) error
	GetBook(book *models.BookResponse) error
	GetAllBooks() ([]models.BookResponse, error)
	DeleteBook(id int) error
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

func (b *BookRepository) AddBook(book *models.BookRequest) error {
	var bookResponse models.BookResponse

	log := utils.GetLogger(b.ctx)

	author, err := b.GetOrCreateAuthor(book.Author.Name)
	if err != nil {
		log.Errorf("Failed to get or create author: %d", err)
		return err
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

	insertQuery := `INSERT INTO user_book (name, date_published, isbn, page_count, user_id, author_id) 
VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = b.DB.DB.Exec(
		insertQuery,
		bookResponse.Name,
		bookResponse.DatePublished,
		bookResponse.ISBN,
		bookResponse.PageCount,
		bookResponse.UserID.ID,
		author.ID,
	)
	if err != nil {
		log.Errorf("Failed to perform an insert query on user book table: %d", err)
		return err
	}

	return nil
}

func (b *BookRepository) UpdateBook(book *models.BookRequest) error {
	log := utils.GetLogger(b.ctx)

	author, err := b.GetOrCreateAuthor(book.Author.Name)
	if err != nil {
		log.Errorf("Failed to get or create author: %d", err)
		return err
	}

	book.Author.ID = author.ID

	updateQuery := "UPDATE user_book SET name=$1, date_published=$2, isbn=$3, page_count=$4, author_id=$5 WHERE id=$6"
	result, err := b.DB.DB.Exec(
		updateQuery, book.Name, book.DatePublished, book.ISBN, book.PageCount, book.Author.ID, book.ID,
	)
	if err != nil {
		log.Errorf("Failed to perform an update query in user book table: %d", err)
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Errorf("Error number of rows affected: %v", err)
		return err
	}

	if affectedRows == 0 {
		log.Errorf("No rows were affected: %d", err)
		return err
	}

	return nil
}

func (b *BookRepository) GetBook(book *models.BookResponse) error {
	var date time.Time

	log := utils.GetLogger(b.ctx)

	getQuery := `
	SELECT b.name, b.date_published, b.isbn, b.page_count, a.name
	FROM user_book AS b
	JOIN author AS a ON b.author_id = a.id
	WHERE b.id = $1
`
	err := b.DB.DB.QueryRow(getQuery, book.ID).Scan(
		&book.Name, &date, &book.ISBN, &book.PageCount, &book.Author.Name,
	)
	if err != nil {
		log.Errorf("Query row on user book table failed: %d", err)
		return err
	}

	book.DatePublished = date.Format("2006-01-02")

	return nil
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

func (b *BookRepository) DeleteBook(id int) error {
	log := utils.GetLogger(b.ctx)

	query := "DELETE FROM user_book WHERE id=$1"
	_, err := b.DB.DB.Exec(query, id)
	if err != nil {
		log.Errorf("Failed to perform delete on a user book table: %d", err)
		return err
	}

	return nil
}

func IsAssigned(log logger.Logger, bookID, userID int, db postgres.DB) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM user_book WHERE id = $1 AND user_id = $2)"

	var exists bool
	err := db.DB.QueryRow(query, bookID, userID).Scan(&exists)
	if err != nil {
		log.Errorf("error checking book assignment: %w", err)
		return false, err
	}

	return exists, nil
}

func ValidateISBNExists(log logger.Logger, isbn string, db postgres.DB) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM user_book WHERE isbn = $1)"

	var exists bool
	err := db.DB.QueryRow(query, isbn).Scan(&exists)
	if err != nil {
		log.Errorf("error validating ISBN %s: %w", isbn, err)
		return false, err
	}

	return exists, nil
}
