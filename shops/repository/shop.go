package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"library/shops/models"
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

func (s *ShopRepository) GetOrCreateAuthor(authorName string) (models.Author, error) {
	var author models.Author

	err := s.DBPool.QueryRow(context.Background(), "SELECT id FROM author WHERE name = $1", authorName).Scan(&author.ID)
	if err != nil {
		err = s.DBPool.QueryRow(context.Background(), "INSERT INTO author (name) VALUES ($1) RETURNING id", authorName).Scan(&author.ID)
		if err != nil {
			errorMessage := "Author error: " + err.Error()
			return author, fmt.Errorf(errorMessage)
		}
	}

	author.Names = authorName

	return author, nil
}

func IsISBN(bookISBN string, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM hobbit WHERE isbn = $1)"

	var exists bool
	err := db.QueryRow(context.Background(), query, bookISBN).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking book assignment: %w", err)
	}

	return exists, nil
}

func (s *ShopRepository) LoadBooks(book models.GoogleBooksRequest) error {
	for _, item := range book.Items {
		var bookID int

		isbn := ConvertIndustryIdentifiers(item.VolumeInfo.ISBN)
		checkISBN, _ := IsISBN(isbn, s.DBPool)
		datePublished, _ := ParseDateString(item.VolumeInfo.DatePublished)

		if !checkISBN {
			err := s.DBPool.QueryRow(context.Background(), "INSERT INTO hobbit (name, date_published, isbn, page_count) VALUES ($1, $2, $3, $4) RETURNING id",
				item.VolumeInfo.Name, datePublished, isbn, item.VolumeInfo.PageCount).Scan(&bookID)
			if err != nil {
				return fmt.Errorf("error inserting into hobbit: %v", err)
			}
		}

		for _, authorName := range item.VolumeInfo.Authors {
			author, err := s.GetOrCreateAuthor(authorName)
			if err != nil {
				return fmt.Errorf("author error: %v", err)
			}

			_, err = s.DBPool.Exec(context.Background(), "INSERT INTO hobbit_authors (book_id, author_id) VALUES ($1, $2)", bookID, author.ID)
			if err != nil {
				return fmt.Errorf("error inserting into authors: %v", err)
			}
		}
	}

	return nil
}
