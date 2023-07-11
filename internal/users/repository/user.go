package repository

import (
	"context"
	"fmt"
	"library/internal/users/models"
	"library/utils"

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
	exists, err := checkUserEmailExist(user.Email, r.DBPool)
	if err != nil {
		errorMessage := "Checking user email error: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	if exists {
		errorMessage := "User email already exists: " + user.Email
		return fmt.Errorf(errorMessage)
	}

	insertQuery := "INSERT INTO users (firstname, lastname, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	result, err := r.DBPool.Exec(context.Background(), insertQuery, user.Firstname, user.Lastname, user.Email, user.Password, user.Role)
	if err != nil {
		errorMessage := "Failed to add user to database: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		errorMessage := "No rows were affected: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func (r *UserRepository) UpdateUser(id int, user *models.User, c *gin.Context) error {
	updateQuery := "UPDATE users SET firstname=$1, lastname=$2, email=$3, role=$4 WHERE id=$5"
	result, err := r.DBPool.Exec(context.Background(), updateQuery, user.Firstname, user.Lastname, user.Email, user.Role, id)
	if err != nil {
		errorMessage := "Failed to execute query: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		errorMessage := "No rows were affected: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func (r *UserRepository) GetUser(id int, user *models.User, c *gin.Context) error {
	getQuery := "SELECT firstname, lastname, email, role FROM users WHERE id=$1"
	err := r.DBPool.QueryRow(context.Background(), getQuery, id).Scan(&user.Firstname, &user.Lastname, &user.Email, &user.Role)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func (r *UserRepository) GetAllUsers(c *gin.Context) ([]models.User, error) {
	var users []models.User

	getAllQuery := "SELECT id, firstname, lastname, email, role FROM users ORDER BY id"
	rows, err := r.DBPool.Query(context.Background(), getAllQuery)
	if err != nil {
		errorMessage := "QueryRow failed: " + err.Error()
		return users, fmt.Errorf(errorMessage)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Role)
		if err != nil {
			errorMessage := "QueryRow failed: " + err.Error()
			return users, fmt.Errorf(errorMessage)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		errorMessage := "Query error: " + err.Error()
		return users, fmt.Errorf(errorMessage)
	}

	return users, nil
}

func (r *UserRepository) DeleteUser(id int, c *gin.Context) error {
	exists, err := utils.CheckIDExists("users", id, r.DBPool)
	if err != nil {
		errorMessage := "Checking user ID error: " + string(rune(id))
		return fmt.Errorf(errorMessage)
	}

	if !exists {
		errorMessage := "User ID doesn't exists: " + string(rune(id))
		return fmt.Errorf(errorMessage)
	}

	query := "DELETE FROM users WHERE id=$1"
	_, err = r.DBPool.Exec(context.Background(), query, id)
	if err != nil {
		errorMessage := fmt.Sprintf("User delete error ID %d: %s", id, err.Error())
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func checkUserEmailExist(email string, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"

	var exists bool
	err := db.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		errorMessage := fmt.Sprintf("Checking user email error %s: %s", email, err.Error())
		return false, fmt.Errorf(errorMessage)
	}

	return exists, nil
}
