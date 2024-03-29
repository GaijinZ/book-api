package pkg

import (
	"library/books/models"
	"library/pkg/logger"
	"time"
)

func CheckPublishedDate(log logger.Logger, book *models.BookRequest) error {
	datePublished, err := time.Parse(time.DateOnly, book.DatePublished)
	if err != nil {
		log.Errorf("Failed to parse date: %s", err)
		return err
	}

	today := time.Now().Local().Truncate(24 * time.Hour)
	if datePublished.After(today) {
		log.Errorf("Date cannot be future date")
		return err
	}

	book.DatePublished = datePublished.Format(time.DateOnly)

	return nil
}
