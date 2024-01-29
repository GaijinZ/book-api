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

	row, err := r.DB.DB.Exec(
		InsertUser,
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
	lastInsertedID, _ := row.LastInsertId()

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

	result, err := r.DB.DB.Exec(
		UpdateUser,
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

	err := r.DB.DB.QueryRow(
		GetUserByID,
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

	rows, err := r.DB.DB.Query(GetUsers)
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

	_, err = r.DB.DB.Exec(DeleteUser, id)
	if err != nil {
		log.Errorf("User delete error ID %d: %s", id, err)
		return 0, err
	}

	return id, nil
}

func checkUserEmailExist(log logger.Logger, email string, db postgres.DB) (bool, error) {
	var exists bool
	err := db.DB.QueryRow(GetUserByEmail, email).Scan(&exists)
	if err != nil {
		log.Errorf("Checking user email error %s: %s", email, err)
		return false, err
	}

	return exists, nil
}
