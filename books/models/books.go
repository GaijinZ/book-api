package models

import (
	"library/users/models"
)

// AuthorRequest represents the request body for creating or updating an author.
type AuthorRequest struct {
	ID   int    `json:"id,omitempty" form:"id"`
	Name string `json:"name,omitempty" form:"name"`
}

// AuthorResponse represents the response body for retrieving author information.
type AuthorResponse struct {
	ID   int    `json:"id,omitempty" form:"id"`
	Name string `json:"name,omitempty" form:"name"`
}

// BookRequest represents the request body for creating or updating a book.
type BookRequest struct {
	ID            int           `json:"id,omitempty" form:"id"`
	Name          string        `json:"name,omitempty" form:"name"`
	DatePublished string        `json:"date_published,omitempty" form:"date_published"`
	ISBN          string        `json:"isbn,omitempty" form:"isbn"`
	PageCount     int           `json:"page_count,omitempty" form:"page_count"`
	UserID        models.User   `json:"user,omitempty"`
	Author        AuthorRequest `json:"author,omitempty"`
}

// BookResponse represents the response body for retrieving book information.
type BookResponse struct {
	ID            int            `json:"id,omitempty" form:"id"`
	Name          string         `json:"name,omitempty" form:"name"`
	DatePublished string         `json:"date_published,omitempty" form:"date_published"`
	ISBN          string         `json:"isbn,omitempty" form:"isbn"`
	PageCount     int            `json:"page_count,omitempty" form:"page_count"`
	UserID        models.User    `json:"user,omitempty"`
	Author        AuthorResponse `json:"author,omitempty"`
}
