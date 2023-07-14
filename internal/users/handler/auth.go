package handler

import (
	"errors"
	"fmt"
	"library/internal/users/models"
	"library/internal/users/repository"
	"library/utils"
	middleware "library/utils/middleware"
	"net/http"
	"time"

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
	fmt.Println("USER", user)
	check := utils.CheckPasswordHash(user.Password, auth.Password)
	if !check {
		err = errors.New("invalid credentials")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		err = errors.New("token generate error")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, int(time.Now().Add(time.Hour*24).Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User logged succesfully"})
}

func (u *UserAuth) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged out"})
}
