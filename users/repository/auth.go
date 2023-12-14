package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"library/pkg/logger"
	"library/users/models"
)

type AutherRepository interface {
	Login(user *models.User, auth *models.Authentication) error
}

type AuthRepository struct {
	ctx    context.Context
	DBPool *pgxpool.Pool
}

func NewAuthRepository(ctx context.Context, dbPool *pgxpool.Pool) AutherRepository {
	return &AuthRepository{
		ctx:    ctx,
		DBPool: dbPool,
	}
}

func (u *AuthRepository) Login(user *models.User, auth *models.Authentication) error {
	log := u.ctx.Value("logger").(logger.Logger)

	query := "SELECT id, firstname, lastname, password, email, role FROM users WHERE email=$1"
	err := u.DBPool.QueryRow(
		context.Background(),
		query,
		auth.Email,
	).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.Role)
	if err != nil {
		log.Errorf("Failed to perform a select query: %v", err)
		return err
	}

	return nil
}
