package models

import (
	"encoding/json"
	"time"
)

type UserTransactionResponse struct {
	ID              int             `json:"id,omitempty"`
	BookList        json.RawMessage `json:"book_list"`
	UserID          int             `json:"user_id,omitempty"`
	Quantity        int             `json:"amount,omitempty"`
	TransactionDate time.Time       `json:"transaction_date"`
}

type UserTransactionRequest struct {
	ID              int       `json:"id,omitempty"`
	BookList        *Book     `json:"book_list"`
	UserID          int       `json:"user_id,omitempty"`
	BookID          int       `json:"book_id,omitempty"`
	Quantity        int       `json:"amount,omitempty"`
	TransactionDate time.Time `json:"transaction_date"`
}

type TransactionResponse struct {
	Quantity int `json:"amount,omitempty"`
}

type Author struct {
	ID   int    `json:"id,omitempty" form:"id"`
	Name string `json:"name,omitempty" form:"name"`
}

type Book struct {
	ID            int       `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	DatePublished time.Time `json:"date_published,omitempty"`
	ISBN          string    `json:"isbn,omitempty"`
	PageCount     int       `json:"page_count,omitempty"`
	Quantity      int       `json:"quantity,omitempty"`
	Author        []Author  `json:"author,omitempty"`
}
