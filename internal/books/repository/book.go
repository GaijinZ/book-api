package repository

import (
	"context"
	"fmt"
	"library/internal/books/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookRepository struct {
	DBPool *pgxpool.Pool
}

func NewBookRepository(dbPool *pgxpool.Pool) *BookRepository {
	return &BookRepository{
		DBPool: dbPool,
	}
}

func (b *BookRepository) GetOrCreateAuthor(authorName string) (models.Author, error) {
	var author models.Author

	err := b.DBPool.QueryRow(context.Background(), "SELECT id FROM author WHERE name = $1", authorName).Scan(&author.ID)
	if err != nil {
		err = b.DBPool.QueryRow(context.Background(), "INSERT INTO author (name) VALUES ($1) RETURNING id", authorName).Scan(&author.ID)
		if err != nil {
			errorMessage := "Author error: " + err.Error()
			return author, fmt.Errorf(errorMessage)
		}
	}

	author.Name = authorName

	return author, nil
}

func (b *BookRepository) AddBook(userID int, book models.Book, c *gin.Context) error {
	author, err := b.GetOrCreateAuthor(book.Author.Name)
	if err != nil {
		errorMessage := "Author error: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	book.Author.ID = author.ID

	_, err = b.DBPool.Exec(context.Background(), "INSERT INTO book (name, date_published, isbn, page_count, user_id, author_id) VALUES ($1, $2, $3, $4, $5, $6)",
		book.Name, book.DatePublished, book.ISBN, book.PageCount, userID, author.ID)
	if err != nil {
		return err
	}
	return nil
}

func (b *BookRepository) UpdateBook(bookID int, book models.Book, c *gin.Context) error {
	author, err := b.GetOrCreateAuthor(book.Author.Name)
	if err != nil {
		errorMessage := "Author error: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	book.Author.ID = author.ID

	updateQuery := "UPDATE book SET name=$1, date_published=$2, isbn=$3, page_count=$4, author_id=$5 WHERE id=$6"
	result, err := b.DBPool.Exec(context.Background(), updateQuery, book.Name, book.DatePublished, book.ISBN, book.PageCount, book.Author.ID, bookID)
	if err != nil {
		errorMessage := "Author error: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		errorMessage := "No rows were affected: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func (b *BookRepository) GetBook(bookID int, book *models.Book, c *gin.Context) error {
	var date time.Time

	getQuery := `
	SELECT b.name, b.date_published, b.isbn, b.page_count, a.name
	FROM book AS b
	JOIN author AS a ON b.author_id = a.id
	WHERE b.id = $1
`
	err := b.DBPool.QueryRow(context.Background(), getQuery, bookID).Scan(
		&book.Name, &date, &book.ISBN, &book.PageCount,
		&book.Author.Name,
	)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	book.DatePublished = date.Format("2006-01-02")

	return nil
}

func (b *BookRepository) GetAllBooks(c *gin.Context) ([]models.Book, error) {
	var books []models.Book
	var date time.Time

	getAllQuery := `
	SELECT b.name, b.date_published, b.isbn, b.page_count, a.name
	FROM book AS b
	JOIN author AS a ON b.author_id = a.id
	ORDER BY a.id
	`
	rows, err := b.DBPool.Query(context.Background(), getAllQuery)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		return books, fmt.Errorf(errorMessage)
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book

		err := rows.Scan(&book.Name, &date, &book.ISBN, &book.PageCount, &book.Author.Name)
		if err != nil {
			errorMessage := "QueryRow failed: " + err.Error()
			return books, fmt.Errorf(errorMessage)
		}

		book.DatePublished = date.Format("2006-01-02")
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		errorMessage := "Query error: " + err.Error()
		return books, fmt.Errorf(errorMessage)
	}

	return books, nil
}

func (b *BookRepository) DeleteBook(id int, c *gin.Context) error {
	query := "DELETE FROM book WHERE id=$1"
	_, err := b.DBPool.Exec(context.Background(), query, id)
	if err != nil {
		errorMessage := fmt.Sprintf("Book delete error ID %d: %s", id, err.Error())
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func IsAssigned(bookID, userID int, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM book WHERE id = $1 AND user_id = $2)"

	var exists bool
	err := db.QueryRow(context.Background(), query, bookID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking book assignment: %w", err)
	}

	return exists, nil
}

func ValidateISBNExists(isbn string, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM book WHERE isbn = $1)"

	var exists bool
	err := db.QueryRow(context.Background(), query, isbn).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error validating ISBN %s: %w", isbn, err)
	}

	return exists, nil
}
