package models

import (
	"encoding/json"
	"errors"
	"regexp"
)

// User represents a user entity.
// swagger:model
type User struct {
	ID        int    `json:"id,omitempty" form:"id"`
	Firstname string `json:"first_name,omitempty" form:"firstname"`
	Lastname  string `json:"last_name,omitempty" form:"lastname"`
	Email     string `json:"email,omitempty" form:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" form:"password"`
	Role      string `json:"role,omitempty" form:"role"`
	IsActive  bool   `json:"is_active,omitempty" form:"is_active"`
}

// UserResponse represents a response for a user entity.
// swagger:model
type UserResponse struct {
	ID        int             `json:"id,omitempty" form:"id"`
	Firstname string          `json:"first_name,omitempty" form:"firstname"`
	Lastname  string          `json:"last_name,omitempty" form:"lastname"`
	Email     string          `json:"email,omitempty" form:"email" validate:"required,email"`
	BookList  json.RawMessage `json:"book_list"`
	Role      string          `json:"role,omitempty" form:"role"`
}

// Authentication represents the authentication credentials.
// swagger:model
type Authentication struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// ValidateUser checks for valid user input in firstname, lastname and password
func (u *User) ValidateUser() error {
	validNameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)

	if !validNameRegex.MatchString(u.Firstname) && len(u.Firstname) < 20 {
		return errors.New("firstname contains forbidden characters")
	}

	if !validNameRegex.MatchString(u.Lastname) && len(u.Firstname) < 20 {
		return errors.New("lastname contains forbidden characters")
	}

	if u.Password == "" {
		return errors.New("password is required for new user creation")
	}

	return nil
}
