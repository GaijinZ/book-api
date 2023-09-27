package pkg

import (
	"fmt"
	"library/books/models"
	"time"
)

func CheckDate(book *models.Book) error {
	datePublished, err := time.Parse(time.DateOnly, book.DatePublished)
	if err != nil {
		errorMessage := "Failed to parse date: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	today := time.Now().Local().Truncate(24 * time.Hour)
	if datePublished.After(today) {
		errorMessage := "date cannot be future date"
		return fmt.Errorf(errorMessage)
	}

	book.DatePublished = datePublished.Format(time.DateOnly)

	return nil
}
