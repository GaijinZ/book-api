package repository

import (
	"context"
	"fmt"
	"net/http"
	"userapi/internal/users/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DBPool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DBPool: dbPool,
	}
}

func (r *UserRepository) AddUser(user *models.User, c *gin.Context) error {
	insertQuery := "INSERT INTO users (firstname, lastname, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	result, err := r.DBPool.Exec(context.Background(), insertQuery, user.Firstname, user.Lastname, user.Email, user.Password, user.Role)
	if err != nil {
		errorMessage := "Failed to add user to database: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return err
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		errorMessage := "No rows were affected: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUser(id int, user *models.User, c *gin.Context) error {
	updateQuery := "UPDATE users SET firstname=$1, lastname=$2, email=$3, role=$4 WHERE id=$5"
	result, err := r.DBPool.Exec(context.Background(), updateQuery, user.Firstname, user.Lastname, user.Email, user.Role, id)
	if err != nil {
		errorMessage := "Failed to update user in database: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return err
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user in database"})
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(id int, user *models.User, c *gin.Context) error {
	getQuery := "SELECT id, firstname, lastname, email, role FROM users WHERE id=$1"
	err := r.DBPool.QueryRow(context.Background(), getQuery, id).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Role)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return err
	}

	return nil
}

func (r *UserRepository) GetAllUsers(c *gin.Context) ([]models.User, error) {
	var users []models.User

	getAllQuery := "SELECT id, firstname, lastname, email, role FROM users ORDER BY id"
	rows, err := r.DBPool.Query(context.Background(), getAllQuery)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return nil, err
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

	return users, nil
}

func (r *UserRepository) DeleteUser(id int, c *gin.Context) error {
	exists, err := checkUserExists(id, r.DBPool)
	if err != nil {
		errorMessage := "Checking user ID error: " + string(rune(id))
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return err
	}

	if !exists {
		errorMessage := "User ID doesn't exists: " + string(rune(id))
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return err
	}

	query := "DELETE FROM users WHERE id=$1"
	_, err = r.DBPool.Exec(context.Background(), query, id)
	if err != nil {
		errorMessage := fmt.Sprintf("User delete error ID %d: %s", id, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return err
	}

	return nil
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
