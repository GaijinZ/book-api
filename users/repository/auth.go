package repository

import (
	"context"
	"library/pkg/postgres"
	"library/pkg/utils"
	"library/users/models"
)

type AutherRepository interface {
	Login(user *models.User, auth *models.Authentication) error
}

type AuthRepository struct {
	ctx context.Context
	db  postgres.DB
}

func NewAuthRepository(ctx context.Context, db postgres.DB) AutherRepository {
	return &AuthRepository{
		ctx: ctx,
		db:  db,
	}
}

func (u *AuthRepository) Login(user *models.User, auth *models.Authentication) error {
	log := utils.GetLogger(u.ctx)

	err := u.db.DB.QueryRow(
		GetUserByEmail,
		auth.Email,
	).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.Role)
	if err != nil {
		log.Errorf("Failed to perform a select query: %v", err)
		return err
	}

	return nil
}
