package handler

import (
	"fmt"
	"library/internal/users/models"
	"library/internal/users/repository"
	"library/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	middleware "library/utils/middleware"
)

type Userer interface {
	AddUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type UserHandler struct {
	userRepository *repository.UserRepository
}

func NewUserHandler(userRepository *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepository: userRepository,
	}
}

var validate = validator.New()

func (h *UserHandler) AddUser(c *gin.Context) {
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

	user.Password, err = utils.GenerateHashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidateUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userRepository.AddUser(&user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User

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

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.userRepository.UpdateUser(userID, &user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"User updated successfully": user})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var user models.User

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

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.userRepository.GetUser(userID, &user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepository.GetAllUsers(c)
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

	deleteUser, err := strconv.Atoi(c.Param("delete_id"))
	if err != nil {
		errorMessage := "Wrong user ID: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	err = h.userRepository.DeleteUser(deleteUser, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
