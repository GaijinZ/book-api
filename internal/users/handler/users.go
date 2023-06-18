package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"userapi/internal/users/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Userer interface {
	AddUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetUser(c *gin.Context)
}

type DBPoolHandler struct {
	DBPool *pgxpool.Pool
}

func (h *DBPoolHandler) AddUser(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertQuery := "INSERT INTO users (firstname, lastname, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	result, err := h.DBPool.Exec(context.Background(), insertQuery, newUser.Firstname, newUser.Lastname, newUser.Email, newUser.Password, newUser.Role)
	if err != nil {
		errorMessage := "Failed to add user to database: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		errorMessage := "No rows were affected: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User added successfully"})
}

func (h *DBPoolHandler) UpdateUser(c *gin.Context) {
	var updateUser models.User

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorMessage := "Wrong user ID: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	if err := c.ShouldBindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQuery := "UPDATE users SET firstname=$1, lastname=$2, email=$3, role=$4 WHERE id=$5"
	result, err := h.DBPool.Exec(context.Background(), updateQuery, updateUser.Firstname, updateUser.Lastname, updateUser.Email, updateUser.Role, userID)
	if err != nil {
		errorMessage := "Failed to update user in database: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user in database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"User updated successfully": updateUser})
}

func (h *DBPoolHandler) GetUser(c *gin.Context) {
	var user models.User

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorMessage := "Wrong user ID: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	getQuery := "SELECT id, firstname, lastname, email, role FROM users WHERE id=$1"
	err = h.DBPool.QueryRow(context.Background(), getQuery, userID).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Role)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *DBPoolHandler) GetAllUsers(c *gin.Context) {
	var users []models.User

	getAllQuery := "SELECT id, firstname, lastname, email, password, role FROM users ORDER BY id"
	rows, err := h.DBPool.Query(context.Background(), getAllQuery)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Role)
		if err != nil {
			errorMessage := "QueryRow failed: " + err.Error()
			c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		errorMessage := "Query error: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *DBPoolHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorMessage := "Wrong user ID: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	exists, err := checkUserExists(userID, h.DBPool)
	if err != nil {
		errorMessage := "Checking user ID error: " + string(rune(userID))
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	if !exists {
		errorMessage := "User ID doesn't exists: " + string(rune(userID))
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	query := "DELETE FROM users WHERE id=$1"
	_, err = h.DBPool.Exec(context.Background(), query, userID)
	if err != nil {
		errorMessage := fmt.Sprintf("User delete error ID %d: %s", userID, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func checkUserExists(userID int, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)"

	var exists bool
	err := db.QueryRow(context.Background(), query, userID).Scan(&exists)
	if err != nil {
		errorMessage := fmt.Sprintf("Checking user ID error %d: %s", userID, err.Error())
		return false, fmt.Errorf(errorMessage)
	}

	return exists, nil
}
