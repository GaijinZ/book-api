package handler

import (
	"context"
	"errors"
	"library/pkg"
	"library/pkg/middleware"
	"library/users/models"
	"library/users/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserAuther interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
}

type UserAuth struct {
	ctx            context.Context
	authRepository repository.AutherRepository
}

func NewUserAuth(ctx context.Context, authRepository repository.AutherRepository) UserAuther {
	return &UserAuth{
		ctx:            ctx,
		authRepository: authRepository,
	}
}

func (u *UserAuth) Login(c *gin.Context) {
	var auth models.Authentication
	var user models.User
	var err error

	if err = c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = u.authRepository.Login(&user, &auth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	check := pkg.CheckPasswordHash(user.Password, auth.Password)
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
	c.JSON(http.StatusOK, gin.H{"message": "User logged successfully"})
}

func (u *UserAuth) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged out"})
}
