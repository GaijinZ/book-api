package models

import (
	"library/users/models"
)

type AuthorRequest struct {
	ID   int    `json:"id,omitempty" form:"id"`
	Name string `json:"name,omitempty" form:"name"`
}

type AuthorResponse struct {
	ID   int    `json:"id,omitempty" form:"id"`
	Name string `json:"name,omitempty" form:"name"`
}

type BookRequest struct {
	ID            int    `json:"id,omitempty" form:"id"`
	Name          string `json:"name,omitempty" form:"name"`
	DatePublished string `json:"date_published,omitempty" form:"date_published"`
	ISBN          string `json:"isbn,omitempty" form:"isbn"`
	PageCount     int    `json:"page_count,omitempty" form:"page_count"`
	UserID        models.User
	Author        AuthorRequest
}

type BookResponse struct {
	ID            int    `json:"id,omitempty" form:"id"`
	Name          string `json:"name,omitempty" form:"name"`
	DatePublished string `json:"date_published,omitempty" form:"date_published"`
	ISBN          string `json:"isbn,omitempty" form:"isbn"`
	PageCount     int    `json:"page_count,omitempty" form:"page_count"`
	UserID        models.User
	Author        AuthorResponse
}
