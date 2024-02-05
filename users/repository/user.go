package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"library/pkg/postgres"
	"library/pkg/rabbitMQ/rabbitMQ"
	"library/pkg/redis"
	"library/pkg/utils"
	"library/users/models"
	"strconv"
	"time"
)

type UsererRepository interface {
	AddUser(user *models.User) (int, error)
	UpdateUser(user *models.User) (*models.UserResponse, error)
	GetUser(id int) (*models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	DeleteUser(id int) (int, error)
	ActivateUser(id int) (int, error)
}

type UserRepository struct {
	ctx         context.Context
	DB          postgres.DB
	redisClient *redis.Client
	rmq         *rabbitMQ.RabbitMQ
}

func NewUserRepository(ctx context.Context, db postgres.DB, redisClient *redis.Client, rmq *rabbitMQ.RabbitMQ) UsererRepository {
	return &UserRepository{
		ctx:         ctx,
		DB:          db,
		redisClient: redisClient,
		rmq:         rmq,
	}
}

func (r *UserRepository) AddUser(user *models.User) (int, error) {
	log := utils.GetLogger(r.ctx)

	exists, err := checkUserEmailExist(user.Email, r.DB.DB)
	if err != nil {
		log.Errorf("Failed to check user email: %v", err)
		return 0, err
	}

	if exists {
		log.Errorf("User email already exists: %v", user.Email)
		return 0, errors.New("user email already exists")
	}

	var lastInsertedID int

	err = r.DB.DB.QueryRow(InsertUser,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&lastInsertedID)

	if err != nil {
		log.Errorf("Failed to add user to database: %v", err)
		return 0, err
	}

	activationLink, err := utils.GenerateActivationLink(r.ctx, r.redisClient, lastInsertedID)
	if err != nil {
		log.Errorf("")
		return 0, err
	}
	r.rmq.Producer(r.ctx, activationLink)

	return lastInsertedID, nil
}

func (r *UserRepository) UpdateUser(user *models.User) (*models.UserResponse, error) {
	var userResponse *models.UserResponse

	log := utils.GetLogger(r.ctx)

	isActivated, err := isUserActivated(r.ctx, user.Email, r.redisClient, r.DB.DB)
	if err != nil {
		log.Errorf("Failed to check if user is activated")
		return userResponse, err
	}

	if !isActivated {
		log.Errorf("User is not activated")
		return userResponse, errors.New("user is not activated")
	}

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
	var err error

	log := utils.GetLogger(r.ctx)

	err = r.DB.DB.QueryRow(
		GetUserByID,
		id,
	).Scan(&userResponse.ID, &userResponse.Firstname, &userResponse.Lastname, &userResponse.Email, &userResponse.Role)
	if err != nil {
		return userResponse, err
	}

	isActivated, err := isUserActivated(r.ctx, userResponse.Email, r.redisClient, r.DB.DB)
	if err != nil {
		log.Errorf("Failed to check if user is activated")
		return userResponse, err
	}

	if !isActivated {
		log.Errorf("User is not activated")
		return userResponse, errors.New("user is not activated")
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

	exists, err := postgres.CheckIDExists("users", id, r.DB.GetDB())
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

func (r *UserRepository) ActivateUser(id int) (int, error) {
	userResponse := &models.UserResponse{}
	var err error

	log := utils.GetLogger(r.ctx)

	_, err = r.redisClient.Client.Get(r.ctx, strconv.Itoa(id)).Result()
	if err != nil {
		fmt.Printf("Error retrieving value from Redis: %v\n", err)
		return 0, err
	}

	userResponse, err = getUserByID(id, userResponse, r.DB.DB)
	if err != nil {
		log.Errorf("QueryRows failed: %v", err)
		return 0, err
	}
	log.Infof("userResponse %v", userResponse)
	if err = activateUser(r.ctx, userResponse, r.redisClient, r.DB.DB); err != nil {
		log.Errorf("Activate user error: %v", err)
		return 0, err
	}

	return id, nil
}

func checkUserEmailExist(email string, db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRow(CheckUserByEmail, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func getUserByID(id int, userResponse *models.UserResponse, db *sql.DB) (*models.UserResponse, error) {
	err := db.QueryRow(
		GetUserByID,
		id,
	).Scan(&userResponse.ID, &userResponse.Firstname, &userResponse.Lastname, &userResponse.Email, &userResponse.Role)
	if err != nil {
		return userResponse, err
	}

	return userResponse, nil
}

func activateUser(ctx context.Context, userResponse *models.UserResponse, redis *redis.Client, db *sql.DB) error {
	_, err := db.Exec(ActivateUser, userResponse.ID)
	if err != nil {
		return err
	}

	_, err = redis.Client.Set(ctx, userResponse.Email, true, 7*24*time.Hour).Result()
	if err != nil {
		return err
	}

	return nil
}

func isUserActivated(ctx context.Context, userEmail string, redis *redis.Client, db *sql.DB) (bool, error) {
	cachedStatus, err := redis.Client.Get(ctx, userEmail).Bool()
	if err == nil {
		return cachedStatus, nil
	}

	var activated bool
	err = db.QueryRow(IsUserActive, userEmail).Scan(&activated)
	if err != nil {
		return false, err
	}

	_, err = redis.Client.Set(ctx, userEmail, activated, 7*24*time.Hour).Result()
	if err != nil {
		return false, err
	}

	return activated, nil
}
