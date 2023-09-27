package utils

import (
	models2 "library/users/models"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(user models2.User) (string, error) {
	claims := &models2.Claims{
		UserID: strconv.Itoa(user.ID),
		Email:  user.Email,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Email,
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (claims *models2.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &models2.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models2.Claims)

	if !ok {
		return nil, err
	}

	return claims, nil
}
