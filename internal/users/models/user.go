package models

import "fmt"

type User struct {
	ID        int    `json:"id,omitempty" form:"id"`
	Firstname string `json:"firstname,omitempty" form:"firstname"`
	Lastname  string `json:"lastname,omitempty" form:"lastname"`
	Email     string `json:"email,omitempty" form:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" form:"password"`
	Role      string `json:"role,omitempty" form:"role"`
}

type Authentication struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (u *User) ValidatePassword() error {
	if u.Password == "" {
		return fmt.Errorf("password is required for new user creation")
	}

	return nil
}
