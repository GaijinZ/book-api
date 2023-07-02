package repository

import (
	"context"
	"library/internal/users/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	DBPool *pgxpool.Pool
}

func NewAuthRepository(dbPool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		DBPool: dbPool,
	}
}

func (u *AuthRepository) Login(user *models.User, auth *models.Authentication, c *gin.Context) error {
	query := "SELECT id, firstname, lastname, password, email, role FROM users WHERE email=$1"
	err := u.DBPool.QueryRow(context.Background(), query, auth.Email).Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Password, &user.Email, &user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	return nil
}
