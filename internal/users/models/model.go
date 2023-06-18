package models

type User struct {
	ID int `json:"id" form:"id"`
	Firstname string `json:"firstname" form:"firstname"`
	Lastname string `json:"lastname" form:"lastname"`
	Email string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Role string `json:"role" form:"role"`
}
