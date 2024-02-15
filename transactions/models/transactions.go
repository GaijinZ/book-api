package models

import (
	"encoding/json"
	"time"
)

// UserTransactionResponse represents a user transaction response.
//	@Summary		User transaction response
//	@Description	Represents a user transaction response
type UserTransactionResponse struct {
	ID              int             `json:"id,omitempty"`
	BookList        json.RawMessage `json:"book_list"`
	UserID          int             `json:"user_id,omitempty"`
	Quantity        int             `json:"amount,omitempty"`
	TransactionDate time.Time       `json:"transaction_date"`
}

// UserTransactionRequest represents a user transaction request.
//	@Summary		User transaction request
//	@Description	Represents a user transaction request
type UserTransactionRequest struct {
	ID              int       `json:"id,omitempty"`
	BookList        *Book     `json:"book_list"`
	UserID          int       `json:"user_id,omitempty"`
	BookID          int       `json:"book_id,omitempty"`
	Quantity        int       `json:"amount,omitempty"`
	TransactionDate time.Time `json:"transaction_date"`
}

// TransactionResponse represents a transaction response.
//	@Summary		Transaction response
//	@Description	Represents a transaction response
type TransactionResponse struct {
	Quantity int `json:"amount,omitempty"`
}

// Author represents an author of a book.
//	@Summary		Author
//	@Description	Represents an author of a book
type Author struct {
	ID   int    `json:"id,omitempty" form:"id"`
	Name string `json:"name,omitempty" form:"name"`
}

// Book represents a book.
//	@Summary		Book
//	@Description	Represents a book
type Book struct {
	ID            int       `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	DatePublished time.Time `json:"date_published,omitempty"`
	ISBN          string    `json:"isbn,omitempty"`
	PageCount     int       `json:"page_count,omitempty"`
	Quantity      int       `json:"quantity,omitempty"`
	Author        []Author  `json:"author,omitempty"`
}
