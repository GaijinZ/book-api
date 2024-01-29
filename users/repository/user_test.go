package repository_test

import (
	"context"
	"errors"
	"fmt"
	"library/pkg/logger"
	"library/pkg/postgres"
	"library/users/models"
	"library/users/repository"

	"github.com/DATA-DOG/go-sqlmock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API Test", func() {
	var (
		newUserRepo  repository.UsererRepository
		user         *models.User
		userResponse *models.UserResponse
		mock         sqlmock.Sqlmock
		fakeDB       *postgres.DB
		err          error
	)

	JustBeforeEach(func() {
		var ctx context.Context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))

		fakeDB, _ = postgres.NewFakeDB(ctx)

		newUserRepo = repository.NewUserRepository(ctx, *fakeDB)

		user = &models.User{
			Firstname: "tmosto",
			Lastname:  "tmosto",
			Email:     "tmosto@elo.com",
			Password:  "hashed_password",
			Role:      "superuser",
		}

		mock = fakeDB.GetMock()
	})

	AfterEach(func() {
		fakeDB.Close()
	})

	Describe("AddUser", func() {
		It("should add a user successfully", func() {
			mock.ExpectQuery(repository.CheckUserByEmail).
				WithArgs(user.Email).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(repository.InsertUser).
				WithArgs(user.Firstname, user.Lastname, user.Email, user.Password, user.Role).
				WillReturnResult(sqlmock.NewResult(1, 1))

			userID, err := newUserRepo.AddUser(user)
			Expect(err).NotTo(HaveOccurred())
			Expect(userID).To(Equal(1), "Expected AddUser to be called with the correct arguments")
		})

		It("should return an error if database query fails", func() {
			mock.ExpectQuery(repository.CheckUserByEmail).
				WithArgs(user.Email).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(repository.InsertUser).
				WithArgs(user.Firstname, user.Lastname, user.Email, user.Password, user.Role).
				WillReturnError(errors.New("database error"))

			_, err = newUserRepo.AddUser(user)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("database error"))
		})

		It("should return an error if user email already exists", func() {
			mock.ExpectQuery(repository.CheckUserByEmail).
				WithArgs(user.Email).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

			_, err = newUserRepo.AddUser(user)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("user email already exists"))
		})

		It("should return an error if arguments do not match", func() {
			mock.ExpectQuery(repository.CheckUserByEmail).
				WithArgs(user.Email).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(repository.InsertUser).
				WithArgs(user.Role).
				WillReturnResult(sqlmock.NewResult(1, 1))

			_, err = newUserRepo.AddUser(user)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("arguments do not match"))
		})
	})

	Describe("UpdateUser", func() {
		It("should update a user successfully", func() {
			rowsAffected := int64(1)
			mock.ExpectExec(repository.UpdateUser).
				WithArgs(user.Firstname, user.Lastname, user.Email, user.Role, user.ID).
				WillReturnResult(sqlmock.NewResult(0, rowsAffected))

			userResponse, err = newUserRepo.UpdateUser(user)
			Expect(err).NotTo(HaveOccurred())
			Expect(userResponse).NotTo(BeNil())
		})

		It("should return an error if no rows were affected", func() {
			rowsAffected := int64(0)
			mock.ExpectExec(repository.UpdateUser).
				WithArgs(user.Firstname, user.Lastname, user.Email, user.Role, user.ID).
				WillReturnResult(sqlmock.NewResult(0, rowsAffected))

			_, err := newUserRepo.UpdateUser(user)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no rows were affected"))
		})

		It("should return an error if database query fails", func() {
			mock.ExpectExec(repository.UpdateUser).
				WithArgs(user.Firstname, user.Lastname, user.Email, user.Role, user.ID).
				WillReturnError(errors.New("database error"))

			_, err := newUserRepo.UpdateUser(user)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("database error"))
		})
	})

	Describe("GetUser", func() {
		It("should get a user", func() {
			expectedUser := &models.UserResponse{
				Firstname: "John",
				Lastname:  "Doe",
				Email:     "john@example.com",
				Role:      "user",
			}

			rows := sqlmock.NewRows([]string{"firstname", "lastname", "email", "role"}).
				AddRow(expectedUser.Firstname, expectedUser.Lastname, expectedUser.Email, expectedUser.Role)
			mock.ExpectQuery(repository.GetUserByID).
				WithArgs(123).
				WillReturnRows(rows)

			userResponse, err = newUserRepo.GetUser(123)
			Expect(err).NotTo(HaveOccurred())
			Expect(userResponse).To(Equal(expectedUser))
		})

		It("should return an error if user does not exist", func() {
			userID := 456

			expectedUser := &models.UserResponse{
				ID:        0,
				Firstname: "",
				Lastname:  "",
				Email:     "",
				BookList:  nil,
				Role:      "",
			}

			rows := sqlmock.NewRows([]string{"firstname", "lastname", "email", "role"}).
				AddRow(expectedUser.Firstname, expectedUser.Lastname, expectedUser.Email, expectedUser.Role)
			mock.ExpectQuery(repository.GetUserByID).
				WithArgs(userID).
				WillReturnRows(rows)

			userResponse, err := newUserRepo.GetUser(userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(userResponse).To(Equal(expectedUser))
		})
	})

	Describe("GetAllUsers", func() {
		It("should return a list of users", func() {
			expectedUsers := []models.UserResponse{
				{Firstname: "John", Lastname: "Doe", Email: "john@example.com", Role: "user"},
				{Firstname: "Jane", Lastname: "Doe", Email: "jane@example.com", Role: "user"},
			}

			rows := sqlmock.NewRows([]string{"firstname", "lastname", "email", "role"})
			for _, user := range expectedUsers {
				rows.AddRow(user.Firstname, user.Lastname, user.Email, user.Role)
			}

			mock.ExpectQuery(repository.GetUsers).
				WillReturnRows(rows)

			actualUsers, err := newUserRepo.GetAllUsers()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualUsers).To(ConsistOf(expectedUsers))
		})

		Context("when there are no users in the database", func() {
			It("should return an empty list", func() {
				mock.ExpectQuery(repository.GetUsers).
					WillReturnRows(sqlmock.NewRows([]string{"firstname", "lastname", "email", "role"}))

				actualUsers, err := newUserRepo.GetAllUsers()
				Expect(err).NotTo(HaveOccurred())
				Expect(actualUsers).To(BeEmpty())
			})
		})

		Context("when an error occurs while fetching users", func() {
			It("should return an error", func() {
				errorMsg := "database error"
				mock.ExpectQuery(repository.GetUsers).
					WillReturnError(errors.New(errorMsg))

				_, err := newUserRepo.GetAllUsers()
				Expect(err).To(MatchError(errorMsg))
			})
		})
	})

	Describe("DeleteUser", func() {
		Context("when user exists", func() {
			It("should delete the user", func() {
				userID := 123

				query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", "users")
				mock.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectExec(repository.DeleteUser).
					WithArgs(userID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				deletedID, err := newUserRepo.DeleteUser(userID)
				Expect(err).NotTo(HaveOccurred())
				Expect(deletedID).To(Equal(userID))
			})
		})

		Context("when user does not exist", func() {
			It("should return an error", func() {
				userID := 123

				query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", "users")
				mock.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

				mock.ExpectExec(repository.DeleteUser).
					WithArgs(userID).
					WillReturnResult(sqlmock.NewResult(0, 0))

				_, err := newUserRepo.DeleteUser(userID)
				Expect(err).To(MatchError("user ID doesn't exists"))
			})
		})

		Context("when an error occurs while deleting user", func() {
			It("should return an error", func() {
				userID := 123
				errorMsg := "database error"

				query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", "users")
				mock.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

				mock.ExpectExec(repository.DeleteUser).
					WithArgs(userID).
					WillReturnError(errors.New(errorMsg))

				_, err := newUserRepo.DeleteUser(userID)
				Expect(err).To(MatchError(errorMsg))
			})
		})
	})
})
