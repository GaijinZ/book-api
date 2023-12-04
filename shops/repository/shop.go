package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"library/shops/models"
	"library/shops/service"
	"time"
)

type ShopRepository struct {
	DBPool *pgxpool.Pool
}

func NewShopRepository(dbPool *pgxpool.Pool) *ShopRepository {
	return &ShopRepository{
		DBPool: dbPool,
	}
}

func (s *ShopRepository) GetOrCreateAuthor(authorName string) (models.Author, error) {
	var author models.Author

	selectAuthor := "SELECT id FROM author WHERE name = $1"
	err := s.DBPool.QueryRow(context.Background(), selectAuthor, authorName).Scan(&author.ID)
	if err != nil {
		insertAuthor := "INSERT INTO author (name) VALUES ($1) RETURNING id"
		err = s.DBPool.QueryRow(context.Background(), insertAuthor, authorName).Scan(&author.ID)
		if err != nil {
			errorMessage := "Author error: " + err.Error()
			return author, fmt.Errorf(errorMessage)
		}
	}

	author.Names = authorName

	return author, nil
}

func (s *ShopRepository) LoadBooks() error {
	booksTypes := []string{"hobbit", "war", "lego"}

	for _, v := range booksTypes {
		bookType, err := searchBooksByTypes(v)
		if err != nil {
		}

		for _, item := range bookType.Items {
			var bookID int

			isbn := ConvertIndustryIdentifiers(item.VolumeInfo.ISBN)
			checkISBN, _ := IsISBN(isbn, s.DBPool)
			datePublished, _ := ParseDateString(item.VolumeInfo.DatePublished)

			if !checkISBN {
				insertBook := "INSERT INTO book (name, date_published, isbn, page_count, quantity)" +
					" VALUES ($1, $2, $3, $4, $5) RETURNING id"
				err = s.DBPool.QueryRow(context.Background(),
					insertBook,
					item.VolumeInfo.Name, datePublished, isbn, item.VolumeInfo.PageCount, 1).Scan(&bookID)
				if err != nil {
					return fmt.Errorf("error inserting into book: %v", err)
				}
			} else {
				updateBook := "UPDATE book SET quantity = quantity + 1 WHERE isbn = $1 RETURNING id"
				err = s.DBPool.QueryRow(context.Background(), updateBook, isbn).Scan(&bookID)
				if err != nil {
					return fmt.Errorf("error updating quantity: %v", err)
				}
			}

			for _, authorName := range item.VolumeInfo.Authors {
				author, err := s.GetOrCreateAuthor(authorName)
				if err != nil {
					return fmt.Errorf("author error: %v", err)
				}

				isAssigned, err := IsAuthorAssigned(author.ID, bookID, s.DBPool)
				if err != nil {
					return fmt.Errorf("isAssigned error: %v", err)
				}

				if !isAssigned {
					insertBookAuthors := "INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)"
					_, err = s.DBPool.Exec(context.Background(), insertBookAuthors, bookID, author.ID)
					if err != nil {
						return fmt.Errorf("error inserting into authors: %v", err)
					}
				}

			}
		}
	}

	return nil
}

func searchBooksByTypes(bookType string) (models.GoogleBooksRequest, error) {
	booksResponse, err := service.GetBooks(bookType)
	if err != nil {
		return booksResponse, fmt.Errorf("error fetching books: %v", err)
	}

	return booksResponse, nil
}

func ConvertIndustryIdentifiers(identifiers []models.IndustryIdentifier) string {
	var result []string

	for _, id := range identifiers {
		result = append(result, id.Identifier)
	}

	return result[0]
}

func ParseDateString(publishedDate string) (string, error) {
	dateFormats := []string{
		"2006-01-02",
		"2006-01",
		"2006",
	}

	for _, format := range dateFormats {
		parsedDate, err := time.Parse(format, publishedDate)
		if err == nil {
			return parsedDate.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("unknown date format: %s", publishedDate)
}

func IsISBN(bookISBN string, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM book WHERE isbn = $1)"

	var exists bool
	err := db.QueryRow(context.Background(), query, bookISBN).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking book assignment: %w", err)
	}

	return exists, nil
}

func IsAuthorAssigned(authorID, bookID int, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM book_authors WHERE book_id = $1 AND author_id = $2)"

	var exists bool
	err := db.QueryRow(context.Background(), query, bookID, authorID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking existence: %v", err)
	}

	return exists, nil
}
