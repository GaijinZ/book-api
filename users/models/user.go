package models

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type User struct {
	ID        int    `json:"id,omitempty" form:"id"`
	Firstname string `json:"first_name,omitempty" form:"firstname"`
	Lastname  string `json:"last_name,omitempty" form:"lastname"`
	Email     string `json:"email,omitempty" form:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" form:"password"`
	Role      string `json:"role,omitempty" form:"role"`
}

type UserResponse struct {
	ID        int             `json:"id,omitempty" form:"id"`
	Firstname string          `json:"first_name,omitempty" form:"firstname"`
	Lastname  string          `json:"last_name,omitempty" form:"lastname"`
	Email     string          `json:"email,omitempty" form:"email" validate:"required,email"`
	BookList  json.RawMessage `json:"book_list"`
	Role      string          `json:"role,omitempty" form:"role"`
}

type Authentication struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type ActivationMessage struct {
	Token  string
	UserID int
	Email  string
}

func (u *User) ValidateUser() error {
	validNameRegex := regexp.MustCompile(`^[a-zA-Z]+$`)

	if !validNameRegex.MatchString(u.Firstname) && len(u.Firstname) < 20 {
		return fmt.Errorf("firstname contains forbidden characters")
	}

	if !validNameRegex.MatchString(u.Lastname) && len(u.Firstname) < 20 {
		return fmt.Errorf("lastname contains forbidden characters")
	}

	if u.Password == "" {
		return fmt.Errorf("password is required for new user creation")
	}

	return nil
}
