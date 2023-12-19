package repository

import (
	"context"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/shops/models"
	"library/shops/service"
	"time"
)

type ShopperRepository interface {
	LoadBooks() error
}

type ShopRepository struct {
	ctx context.Context
	DB  postgres.DB
}

func NewShopRepository(ctx context.Context, db postgres.DB) ShopperRepository {
	return &ShopRepository{
		ctx: ctx,
		DB:  db,
	}
}

func (s *ShopRepository) GetOrCreateAuthor(authorName string) (models.Author, error) {
	var author models.Author
	log := s.ctx.Value("logger").(logger.Logger)

	selectAuthor := "SELECT id FROM author WHERE name = $1"
	err := s.DB.DB.QueryRow(selectAuthor, authorName).Scan(&author.ID)
	if err != nil {
		insertAuthor := "INSERT INTO author (name) VALUES ($1) RETURNING id"
		err = s.DB.DB.QueryRow(insertAuthor, authorName).Scan(&author.ID)
		if err != nil {
			log.Errorf("Failed to perform an insert query on author table: %v", err)
			return author, err
		}
	}

	author.Names = authorName

	return author, nil
}

func (s *ShopRepository) LoadBooks() error {
	booksTypes := []string{"hobbit", "war", "lego"}
	log := s.ctx.Value("logger").(logger.Logger)

	for _, v := range booksTypes {
		bookType, err := searchBooksByTypes(log, v)
		if err != nil {
		}

		for _, item := range bookType.Items {
			var bookID int

			isbn := ConvertIndustryIdentifiers(item.VolumeInfo.ISBN)
			checkISBN, _ := IsISBN(log, isbn, s.DB)
			datePublished, _ := ParseDateString(log, item.VolumeInfo.DatePublished)

			if !checkISBN {
				insertBook := "INSERT INTO book (name, date_published, isbn, page_count, quantity)" +
					" VALUES ($1, $2, $3, $4, $5) RETURNING id"
				err = s.DB.DB.QueryRow(
					insertBook,
					item.VolumeInfo.Name, datePublished, isbn, item.VolumeInfo.PageCount, 1).Scan(&bookID)
				if err != nil {
					log.Errorf("Failed to perform an insert query on book table: %v", err)
					return err
				}
			} else {
				updateBook := "UPDATE book SET quantity = quantity + 1 WHERE isbn = $1 RETURNING id"
				err = s.DB.DB.QueryRow(updateBook, isbn).Scan(&bookID)
				if err != nil {
					log.Errorf("Failed to perform an update query on book table: %v", err)
					return err
				}
			}

			for _, authorName := range item.VolumeInfo.Authors {
				author, err := s.GetOrCreateAuthor(authorName)
				if err != nil {
					log.Errorf("Failed to get or create author: %v", err)
					return err
				}

				isAssigned, err := IsAuthorAssigned(log, author.ID, bookID, s.DB)
				if err != nil {
					log.Errorf("Failed to check if book is already assigned: %v", err)
					return err
				}

				if !isAssigned {
					insertBookAuthors := "INSERT INTO book_authors (book_id, author_id) VALUES ($1, $2)"
					_, err = s.DB.DB.Exec(insertBookAuthors, bookID, author.ID)
					if err != nil {
						log.Errorf("Failed to perform an insert query on book table: %v", err)
						return err
					}
				}

			}
		}
	}

	return nil
}

func searchBooksByTypes(log logger.Logger, bookType string) (models.BooksRequest, error) {
	booksResponse, err := service.GetBooks(bookType)
	if err != nil {
		log.Errorf("Failed to fetch books: %v", err)
		return booksResponse, err
	}

	return booksResponse, nil
}

func ConvertIndustryIdentifiers(identifiers []models.ISBN) string {
	var result []string

	for _, id := range identifiers {
		result = append(result, id.Identifier)
	}

	return result[0]
}

func ParseDateString(log logger.Logger, publishedDate string) (string, error) {
	var err error
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

	log.Errorf("Unknown date format: %s", publishedDate)
	return "", err
}

func IsISBN(log logger.Logger, bookISBN string, db postgres.DB) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM book WHERE isbn = $1)"

	var exists bool
	err := db.DB.QueryRow(query, bookISBN).Scan(&exists)
	if err != nil {
		log.Errorf("Failed to check book ISBN: %v", err)
		return false, err
	}

	return exists, nil
}

func IsAuthorAssigned(log logger.Logger, authorID, bookID int, db postgres.DB) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM book_authors WHERE book_id = $1 AND author_id = $2)"

	var exists bool
	err := db.DB.QueryRow(query, bookID, authorID).Scan(&exists)
	if err != nil {
		log.Errorf("Failed to check if author is already assigned: %v", err)
		return false, err
	}

	return exists, nil
}
