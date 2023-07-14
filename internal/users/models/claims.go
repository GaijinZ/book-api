package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserID      string `json:"id"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
	Role        string `json:"role"`
	jwt.StandardClaims
}
