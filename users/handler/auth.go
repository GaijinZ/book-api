package handler

import (
	"context"
	"errors"
	"library/pkg"
	"library/pkg/middleware"
	"library/pkg/utils"
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

// Login authenticates a user and generates a JWT token.
//
//	@Summary		Authenticate user
//	@Description	Authenticate user and generate JWT token
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.Authentication							true	"User credentials"
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/v1/users/login [post]
func (u *UserAuth) Login(c *gin.Context) {
	var auth models.Authentication
	var user models.User
	var err error

	log := utils.GetLogger(u.ctx)

	if err = c.ShouldBindJSON(&auth); err != nil {
		log.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = u.authRepository.Login(&user, &auth)
	if err != nil {
		log.Errorf("Error repository login: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	check := pkg.CheckPasswordHash(user.Password, auth.Password)
	if !check {
		log.Errorf("invalid credentials")
		err = errors.New("invalid credentials")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(user)
	if err != nil {
		log.Errorf("token generate error: %v", err)
		err = errors.New("token generate error")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, int(time.Now().Add(time.Hour*24).Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "User logged successfully"})
}

// Logout revokes the user's JWT token.
//
//	@Summary		Logout user
//	@Description	Revoke user's JWT token
//	@Produce		json
//	@Success		200
//	@Router			/v1/users/logout [post]
func (u *UserAuth) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged out"})
}
