package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"library/pkg"
	"library/users/models"
	"library/users/repository"
	"log"
	"net/http"

	"library/pkg/middleware"
)

type Userer interface {
	AddUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type UserHandler struct {
	ctx            context.Context
	userRepository repository.UsererRepository
}

func NewUserHandler(ctx context.Context, userRepository repository.UsererRepository) Userer {
	return &UserHandler{
		ctx:            ctx,
		userRepository: userRepository,
	}
}

func (h *UserHandler) AddUser(c *gin.Context) {
	var validate = validator.New()
	var user models.User
	var err error

	if err = c.ShouldBindJSON(&user); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user.Password, err = pkg.GenerateHashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidateUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userRepository.AddUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = c.GetInt("userID")

	err := h.userRepository.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"User updated successfully": user})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var user models.UserResponse

	user.ID = c.GetInt("userID")

	err := h.userRepository.GetUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := middleware.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if claims.Role != "superuser" {
		errorMessage := fmt.Errorf("not enough permissions")
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
		return
	}

	deleteUser := c.GetInt("delete_id")

	err = h.userRepository.DeleteUser(deleteUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
