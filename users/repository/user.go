package repository

import (
	"context"
	"errors"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/pkg/utils"
	"library/users/models"
)

type UsererRepository interface {
	AddUser(user *models.User) (int, error)
	UpdateUser(user *models.User) (*models.UserResponse, error)
	GetUser(id int) (*models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	DeleteUser(id int) (int, error)
	GetDBPool() postgres.DB
}

type UserRepository struct {
	ctx context.Context
	DB  postgres.DB
}

func NewUserRepository(ctx context.Context, db postgres.DB) UsererRepository {
	return &UserRepository{
		ctx: ctx,
		DB:  db,
	}
}

func (r *UserRepository) GetDBPool() postgres.DB {
	return r.DB
}

func (r *UserRepository) AddUser(user *models.User) (int, error) {
	log := utils.GetLogger(r.ctx)

	exists, err := checkUserEmailExist(log, user.Email, r.DB)
	if err != nil {
		log.Errorf("Failed to check user email: %v", err)
		return 0, err
	}

	if exists {
		log.Errorf("User email already exists: %v", user.Email)
		return 0, errors.New("user email already exists")
	}

	insertQuery := "INSERT INTO users (firstname, lastname, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	result, err := r.DB.DB.Exec(
		insertQuery,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Password,
		user.Role,
	)
	if err != nil {
		log.Errorf("Failed to add user to database: %v", err)
		return 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		log.Errorf("Failed to get last inserted ID: %v", err)
		return 0, errors.New("error getting last insert ID")
	}

	return int(lastInsertedID), nil
}

func (r *UserRepository) UpdateUser(user *models.User) (*models.UserResponse, error) {
	var userResponse *models.UserResponse

	log := utils.GetLogger(r.ctx)

	userResponse = &models.UserResponse{
		ID:        user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Role:      user.Role,
	}

	updateQuery := "UPDATE users SET firstname=$1, lastname=$2, email=$3, role=$4 WHERE id=$5"
	result, err := r.DB.DB.Exec(
		updateQuery,
		userResponse.Firstname,
		userResponse.Lastname,
		userResponse.Email,
		userResponse.Role,
		userResponse.ID,
	)
	if err != nil {
		log.Errorf("Failed to execute query: %v", err)
		return userResponse, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		log.Errorf("Error number of rows affected: %v", err)
		return userResponse, err
	}

	if affectedRows == 0 {
		log.Errorf("No rows were affected: %v", err)
		return userResponse, errors.New("no rows were affected")
	}

	return userResponse, nil
}

func (r *UserRepository) GetUser(id int) (*models.UserResponse, error) {
	userResponse := &models.UserResponse{}

	log := utils.GetLogger(r.ctx)

	getQuery := "SELECT firstname, lastname, email, role FROM users WHERE id=$1"
	err := r.DB.DB.QueryRow(
		getQuery,
		id,
	).Scan(&userResponse.Firstname, &userResponse.Lastname, &userResponse.Email, &userResponse.Role)
	if err != nil {
		log.Errorf("QueryRows failed: %v", err)
		return userResponse, err
	}

	return userResponse, nil
}

func (r *UserRepository) GetAllUsers() ([]models.UserResponse, error) {
	var users []models.UserResponse
	log := utils.GetLogger(r.ctx)

	getAllQuery := "SELECT firstname, lastname, email, role FROM users ORDER BY id"
	rows, err := r.DB.DB.Query(getAllQuery)
	if err != nil {
		log.Errorf("QueryRows failed: %v", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserResponse

		err = rows.Scan(&user.Firstname, &user.Lastname, &user.Email, &user.Role)
		if err != nil {
			log.Errorf("QueryRows failed: %v", err)
			return users, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Errorf("Query failed: %v", err)
		return users, err
	}

	return users, nil
}

func (r *UserRepository) DeleteUser(id int) (int, error) {
	log := utils.GetLogger(r.ctx)

	exists, err := postgres.CheckIDExists("users", id, r.DB)
	if err != nil {
		log.Errorf("Checking user ID error: %v", err)
		return 0, err
	}

	if !exists {
		log.Errorf("User ID doesn't exists: %v", id)
		return 0, errors.New("user ID doesn't exists")
	}

	query := "DELETE FROM users WHERE id=$1"
	_, err = r.DB.DB.Exec(query, id)
	if err != nil {
		log.Errorf("User delete error ID %d: %s", id, err)
		return 0, err
	}

	return id, nil
}

func checkUserEmailExist(log logger.Logger, email string, db postgres.DB) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"

	var exists bool
	err := db.DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		log.Errorf("Checking user email error %s: %s", email, err)
		return false, err
	}

	return exists, nil
}
