package service

import (
	"encoding/json"
	"fmt"
	"library/shops/models"
	"net/http"
)

func GetBooks(bookTitle string) (models.GoogleBooksRequest, error) {
	var googleBooksRequest models.GoogleBooksRequest

	url := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s", bookTitle)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&googleBooksRequest)
	if err != nil {
		return models.GoogleBooksRequest{}, err
	}

	return googleBooksRequest, nil
}
