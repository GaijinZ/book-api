package repository

import (
	"context"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/users/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsererRepository interface {
	AddUser(user *models.User) error
	UpdateUser(id int, user *models.User) error
	GetUser(id int, user *models.User) error
	GetAllUsers() ([]models.User, error)
	DeleteUser(id int) error
}

type UserRepository struct {
	ctx    context.Context
	DBPool *pgxpool.Pool
}

func NewUserRepository(ctx context.Context, dbPool *pgxpool.Pool) UsererRepository {
	return &UserRepository{
		ctx:    ctx,
		DBPool: dbPool,
	}
}

func (r *UserRepository) AddUser(user *models.User) error {
	log := r.ctx.Value("logger").(logger.Logger)

	exists, err := checkUserEmailExist(log, user.Email, r.DBPool)
	if err != nil {
		log.Errorf("Failed to check user email: %v", err)
		return err
	}

	if !exists {
		insertQuery := "INSERT INTO users (firstname, lastname, email, password, role) VALUES ($1, $2, $3, $4, $5)"
		result, err := r.DBPool.Exec(
			context.Background(),
			insertQuery,
			user.Firstname,
			user.Lastname,
			user.Email,
			user.Password,
			user.Role,
		)
		if err != nil {
			log.Errorf("Failed to add user to database: %v", err)
			return err
		}

		affectedRows := result.RowsAffected()

		if affectedRows == 0 {
			log.Errorf("No rows were affected: %v", err)
			return err
		}
	}

	log.Errorf("User email already exists: %v", user.Email)
	return nil
}

func (r *UserRepository) UpdateUser(id int, user *models.User) error {
	log := r.ctx.Value("logger").(logger.Logger)

	updateQuery := "UPDATE users SET firstname=$1, lastname=$2, email=$3, role=$4 WHERE id=$5"
	result, err := r.DBPool.Exec(
		context.Background(),
		updateQuery,
		user.Firstname,
		user.Lastname,
		user.Email,
		user.Role,
		id,
	)
	if err != nil {
		log.Errorf("Failed to execute query: %v", err)
		return err
	}

	affectedRows := result.RowsAffected()

	if affectedRows == 0 {
		log.Errorf("No rows were affected: %v", err)
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(id int, user *models.User) error {
	log := r.ctx.Value("logger").(logger.Logger)

	getQuery := "SELECT firstname, lastname, email, role FROM users WHERE id=$1"
	err := r.DBPool.QueryRow(
		context.Background(),
		getQuery,
		id,
	).Scan(&user.Firstname, &user.Lastname, &user.Email, &user.Role)
	if err != nil {
		log.Errorf("QueryRows failed: %v", err)
		return err
	}

	return nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	log := r.ctx.Value("logger").(logger.Logger)

	getAllQuery := "SELECT firstname, lastname, email, role FROM users ORDER BY id"
	rows, err := r.DBPool.Query(context.Background(), getAllQuery)
	if err != nil {
		log.Errorf("QueryRows failed: %v", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User

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

func (r *UserRepository) DeleteUser(id int) error {
	log := r.ctx.Value("logger").(logger.Logger)

	exists, err := postgres.CheckIDExists("users", id, r.DBPool)
	if err != nil {
		log.Errorf("Checking user ID error: %v", err)
		return err
	}

	if !exists {
		log.Errorf("User ID doesn't exists: %v", string(rune(id)))
		return err
	}

	query := "DELETE FROM users WHERE id=$1"
	_, err = r.DBPool.Exec(context.Background(), query, id)
	if err != nil {
		log.Errorf("User delete error ID %d: %s", id, err)
		return err
	}

	return nil
}

func checkUserEmailExist(log logger.Logger, email string, db *pgxpool.Pool) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"

	var exists bool
	err := db.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		log.Errorf("Checking user email error %s: %s", email, err)
		return false, err
	}

	return exists, nil
}
