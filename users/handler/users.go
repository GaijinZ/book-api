package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"library/pkg"
	"library/pkg/middleware"
	"library/pkg/utils"
	"library/users/models"
	"library/users/repository"
	"net/http"
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

	log := utils.GetLogger(h.ctx)

	if err = c.ShouldBindJSON(&user); err != nil {
		log.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = validate.Struct(user); err != nil {
		log.Warningf("Validation error: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user.Password, err = pkg.GenerateHashPassword(user.Password)
	if err != nil {
		log.Errorf("Error generating hashed password: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidateUser()
	if err != nil {
		log.Warningf("User validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.userRepository.AddUser(&user)
	if err != nil {
		log.Errorf("Error adding user to the repository: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("User added successfully: %v", userID)
	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	var userResponse *models.UserResponse

	log := utils.GetLogger(h.ctx)

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = c.GetInt("userID")

	userResponse, err := h.userRepository.UpdateUser(&user)
	if err != nil {
		log.Errorf("Error updating user in the repository: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("User updated successfully: %v", userResponse.ID)
	c.JSON(http.StatusCreated, gin.H{"User updated successfully": userResponse})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	var user *models.UserResponse
	var err error

	log := utils.GetLogger(h.ctx)

	userID := c.GetInt("userID")

	user, err = h.userRepository.GetUser(userID)
	if err != nil {
		log.Errorf("Error getting user from the repository: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("User: %v", user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	log := utils.GetLogger(h.ctx)

	users, err := h.userRepository.GetAllUsers()
	if err != nil {
		log.Errorf("Error getting all users from the repository: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Users: %v", users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	log := utils.GetLogger(h.ctx)

	token, err := c.Cookie("token")
	if err != nil {
		log.Errorf("Error getting token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := middleware.VerifyJWT(token)
	if err != nil {
		log.Errorf("Authorization failed : %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if claims.Role != "superuser" {
		log.Warningf("not enough permissions")
		errorMessage := fmt.Errorf("not enough permissions")
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
		return
	}

	id := c.GetInt("deleteID")

	deletedID, err := h.userRepository.DeleteUser(id)
	if err != nil {
		log.Warningf("Error deleting user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("User deleted successfully, id: %v", deletedID)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
