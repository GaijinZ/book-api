package repository

import (
	"context"
	"library/books/models"
	"library/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookerRepository interface {
	AddBook(userID int, book models.Book) error
	UpdateBook(bookID int, book models.Book) error
	GetBook(bookID int, book *models.Book) error
	GetAllBooks() ([]models.Book, error)
	DeleteBook(id int) error
	GetDBPool() *pgxpool.Pool
}

type BookRepository struct {
	ctx    context.Context
	DBPool *pgxpool.Pool
}

func NewBookRepository(ctx context.Context, dbPool *pgxpool.Pool) BookerRepository {
	return &BookRepository{
		ctx:    ctx,
		DBPool: dbPool,
	}
}

func (b *BookRepository) GetDBPool() *pgxpool.Pool {
	return b.DBPool
}

func (b *BookRepository) GetOrCreateAuthor(authorName string) (models.Author, error) {
	var author models.Author
	log := b.ctx.Value("logger").(logger.Logger)

	query := "SELECT id FROM author WHERE name = $1"
	err := b.DBPool.QueryRow(b.ctx, query, authorName).Scan(&author.ID)
	if err != nil {
		query = "INSERT INTO author (name) VALUES ($1) RETURNING id"
		err = b.DBPool.QueryRow(context.Background(), query, authorName).Scan(&author.ID)
		if err != nil {
			log.Errorf("Failed to perform an insert query on author table: %d", err)
			return author, err
		}
	}

	author.Name = authorName

	return author, nil
}

func (b *BookRepository) AddBook(userID int, book models.Book) error {
	log := b.ctx.Value("logger").(logger.Logger)

	author, err := b.GetOrCreateAuthor(book.Author.Name)
	if err != nil {
		log.Errorf("Failed to get or create author: %d", err)
		return err
	}

	book.Author.ID = author.ID

	insertQuery := `INSERT INTO user_book (name, date_published, isbn, page_count, user_id, author_id) 
VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = b.DBPool.Exec(
		context.Background(),
		insertQuery,
		book.Name, book.DatePublished, book.ISBN, book.PageCount, userID, author.ID,
	)
	if err != nil {
		log.Errorf("Failed to perform an insert query on user book table: %d", err)
		return err
	}
	return nil
}

func (b *BookRepository) UpdateBook(bookID int, book models.Book) error {
	log := b.ctx.Value("logger").(logger.Logger)

	author, err := b.GetOrCreateAuthor(book.Author.Name)
	if err != nil {
		log.Errorf("Failed to get or create author: %d", err)
		return err
	}

	book.Author.ID = author.ID

	updateQuery := "UPDATE user_book SET name=$1, date_published=$2, isbn=$3, page_count=$4, author_id=$5 WHERE id=$6"
	result, err := b.DBPool.Exec(
		context.Background(),
		updateQuery, book.Name, book.DatePublished, book.ISBN, book.PageCount, book.Author.ID, bookID,
	)
	if err != nil {
		log.Errorf("Failed to perform an update query in user book table: %d", err)
		return err
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		log.Errorf("No rows were affected: %d", err)
		return err
	}

	return nil
}

func (b *BookRepository) GetBook(bookID int, book *models.Book) error {
	var date time.Time
	log := b.ctx.Value("logger").(logger.Logger)

	getQuery := `
	SELECT b.name, b.date_published, b.isbn, b.page_count, a.name
	FROM user_book AS b
	JOIN author AS a ON b.author_id = a.id
	WHERE b.id = $1
`
	err := b.DBPool.QueryRow(context.Background(), getQuery, bookID).Scan(
		&book.Name, &date, &book.ISBN, &book.PageCount,
		&book.Author.Name,
	)
	if err != nil {
		log.Errorf("Query row on user book table failed: %d", err)
		return err
	}

	book.DatePublished = date.Format("2006-01-02")

	return nil
}

func (b *BookRepository) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	var date time.Time
	log := b.ctx.Value("logger").(logger.Logger)

	getAllQuery := `
	SELECT b.name, b.date_published, b.isbn, b.page_count, a.name
	FROM user_book AS b
	JOIN author AS a ON b.author_id = a.id
	ORDER BY a.id
	`
	rows, err := b.DBPool.Query(context.Background(), getAllQuery)
	if err != nil {
		log.Errorf("Failed to perform a query on user book table: %d", err)
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book

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
	log := b.ctx.Value("logger").(logger.Logger)

	query := "DELETE FROM user_book WHERE id=$1"
	_, err := b.DBPool.Exec(context.Background(), query, id)
	if err != nil {
		log.Errorf("Failed to perform delete on a user book table: %d", err)
		return err
	}

	return nil
}

// don't pass db as a value
func IsAssigned(log logger.Logger, bookID, userID int, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM user_book WHERE id = $1 AND user_id = $2)"

	var exists bool
	err := db.QueryRow(context.Background(), query, bookID, userID).Scan(&exists)
	if err != nil {
		log.Errorf("error checking book assignment: %w", err)
		return false, err
	}

	return exists, nil
}

// don't pass db as a value
func ValidateISBNExists(log logger.Logger, isbn string, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM user_book WHERE isbn = $1)"

	var exists bool
	err := db.QueryRow(context.Background(), query, isbn).Scan(&exists)
	if err != nil {
		log.Errorf("error validating ISBN %s: %w", isbn, err)
		return false, err
	}

	return exists, nil
}
