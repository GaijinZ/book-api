package handler

import (
	"errors"
	"net/http"
	"userapi/internal/users/models"
	"userapi/internal/users/repository"
	"userapi/util"

	"github.com/gin-gonic/gin"
)

type UserAuth struct {
	authRepository *repository.AuthRepository
}

func NewUserAuth(authRepository *repository.AuthRepository) *UserAuth {
	return &UserAuth{
		authRepository: authRepository,
	}
}

func (u *UserAuth) Login(c *gin.Context) {
	var auth models.Authentication
	var user models.User
	var err error

	if err = c.Bind(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = u.authRepository.Login(&user, &auth, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	check := util.CheckPasswordHash(user.Password, auth.Password)
	if !check {
		err = errors.New("invalid credentials")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged succesfully"})
}
