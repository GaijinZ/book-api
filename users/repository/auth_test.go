package repository_test

import (
	"context"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/users/models"
	"library/users/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestAuthRepository_Login(t *testing.T) {
	var ctx context.Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))

	fakeDB, _ := postgres.NewFakeDB(ctx)

	newAuthRepo := repository.NewAuthRepository(ctx, *fakeDB)
	testUser := &models.User{}
	testAuth := &models.Authentication{Email: "test@example.com"}

	mock := fakeDB.GetMock()

	mock.ExpectQuery(repository.GetUserByEmail).
		WithArgs(testAuth.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "password", "email", "role"}).
			AddRow(1, "John", "Doe", "hashed_password", testAuth.Email, "user"))

	err := newAuthRepo.Login(testUser, testAuth)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedUser := &models.User{
		ID:        1,
		Firstname: "John",
		Lastname:  "Doe",
		Password:  "hashed_password",
		Email:     testAuth.Email,
		Role:      "user",
	}
	if *testUser != *expectedUser {
		t.Errorf("expected user %+v, got %+v", *expectedUser, *testUser)
	}
}
