package handler

import (
	"context"
	"library/pkg"
	"library/pkg/utils"
	"library/users/models"
	"library/users/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Userer interface {
	AddUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	DeleteUser(c *gin.Context)
	ActivateAccount(c *gin.Context)
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

// AddUser creates a new user.
//
//	@Summary		Create a new user
//	@Description	Creates a new user with the provided data
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User										true	"User object to be created"
//	@Success		201
//	@Failure		400
//	@Failure		422
//	@Failure		500
//	@Router			/v1/users [post]
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

// UpdateUser updates a user data.
//
//	@Summary		Updates a user data
//	@Description	Updates a user with the provided data
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int												true	"User ID"
//	@Param			user	body		models.User											true	"Updated user object"
//	@Success		201
//	@Failure		400
//	@Failure		500
//	@Router			/users/{userID} [put]
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

// GetUser retrieves user details by ID.
//
//	@Summary		Get user details
//	@Description	Retrieves user details by the provided ID
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			user_id	path		int										true	"User ID"
//	@Success		200
//	@Failure		404
//	@Failure		500
//	@Router			/v1/users/{user_id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	var user *models.UserResponse
	var err error

	log := utils.GetLogger(h.ctx)

	getUser, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user, err = h.userRepository.GetUser(getUser)
	if err != nil {
		log.Errorf("Error getting user from the repository: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("User: %v", user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetAllUsers retrieves all users.
//
//	@Summary		Get all users
//	@Description	Retrieves all users
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200
//	@Failure		400
//	@Router			/v1/users [get]
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

// DeleteUser deletes a user.
//
//	@Summary		Delete a user
//	@Description	Deletes a user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			deleteID	path		int												true	"User ID to delete"
//	@Param			delete_id	path		int												true	"Delete ID to authorize"
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Router			/v1/users/{user_id}/{delete_id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	log := utils.GetLogger(h.ctx)

	id := c.GetInt("deleteID")
	claims, exists := c.Get("claims")
	if !exists {
		log.Errorf("Could not get claims from context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not get claims from context"})
		return
	}

	claimsData, ok := claims.(*models.Claims)
	if !ok {
		log.Errorf("Failed to convert claims")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to convert claims"})
		return
	}

	if claimsData.Role != "superuser" {
		log.Warningf("not enough permissions")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not enough permissions"})
		return
	}

	deletedID, err := h.userRepository.DeleteUser(id)
	if err != nil {
		log.Warningf("Error deleting user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("User deleted successfully, id: %v", deletedID)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ActivateAccount activates a user account.
//
//	@Summary		Activate a user account
//	@Description	Activates a user account
//	@Accept			json
//	@Produce		json
//	@Param			userID	query		int											true	"User ID to activate"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/v1/users/activate [get]
func (h *UserHandler) ActivateAccount(c *gin.Context) {
	log := utils.GetLogger(h.ctx)
	userIDStr := c.Query("userID")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Errorf("Invalid user ID: %v", userIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	id, err := h.userRepository.ActivateUser(userID)
	if err != nil {
		log.Errorf("Failed to activate user, id: %v", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}

	log.Infof("User account has been activated: %v", id)
	c.JSON(http.StatusOK, gin.H{"message": "User account activated"})
}
