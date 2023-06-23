package models

type User struct {
	ID        int    `json:"id" form:"id"`
	Firstname string `json:"firstname,omitempty" form:"firstname"`
	Lastname  string `json:"lastname,omitempty" form:"lastname"`
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" form:"password" binding:"required"`
	Role      string `json:"role,omitempty" form:"role"`
}

type Authentication struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
