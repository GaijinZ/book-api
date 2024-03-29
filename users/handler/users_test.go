package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"io"
	"library/pkg/config"
	"library/pkg/middleware"
	"library/pkg/rabbitMQ/rabbitMQ"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"library/pkg/logger"
	"library/users/handler"
	"library/users/models"
	"library/users/repository/repositoryfakes"
	"library/users/server"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API Test", func() {
	var (
		cfg             config.GlobalEnv
		w               *httptest.ResponseRecorder
		ginCtx          *gin.Context
		fakeAuther      *repositoryfakes.FakeAutherRepository
		fakeUserer      *repositoryfakes.FakeUsererRepository
		fakeAuthUser    handler.UserAuther
		fakeUserHandler handler.Userer
		router          *gin.Engine
		request         *models.User
		enricher        func(req *http.Request) *http.Request
		cookie          *http.Cookie
	)

	JustBeforeEach(func() {
		w = httptest.NewRecorder()
		ginCtx, _ = gin.CreateTestContext(w)

		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))
		if err := envconfig.Process("bookapi", &cfg); err != nil {
			log.Fatalf(err.Error())
		}

		rmq, err := rabbitMQ.NewConn(cfg)
		if err != nil {
			log.Fatalf("Failed to create RabbitMQ instance: %v", err)
		}
		defer rmq.Close()

		fakeAuther = &repositoryfakes.FakeAutherRepository{}
		fakeUserer = &repositoryfakes.FakeUsererRepository{}

		fakeAuthUser = handler.NewUserAuth(ctx, fakeAuther)
		fakeUserHandler = handler.NewUserHandler(ctx, fakeUserer)

		router = server.NewRouter(fakeAuthUser, fakeUserHandler)

		request = &models.User{
			ID:        1,
			Firstname: "tmostowashere",
			Lastname:  "tmostowashere",
			Email:     "tmostowashere@tmostowashere.com",
			Password:  "tmostowashere",
			Role:      "superuser",
		}

		token, _ := middleware.GenerateJWT(*request)

		cookie = &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			Domain:   "localhost",
			Expires:  time.Now().Add(time.Hour),
			Secure:   false,
			HttpOnly: true,
		}

		enricher = func(req *http.Request) *http.Request {
			req.Header.Set("Authorization", token)
			req.Header.Set("Content-Type", "application/json")
			return req
		}
	})

	Describe("AddUser", func() {
		It("should add a user successfully", func() {
			body, err := json.Marshal(request)
			Expect(err).To(BeNil())

			ginCtx.Request, err = http.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeUserer.AddUserReturns(1, nil)

			router.ServeHTTP(w, ginCtx.Request)

			args := fakeUserer.AddUserArgsForCall(0)
			Expect(args.ID).To(Equal(1), "Expected AddUser to be called with the correct arguments")
		})
	})

	Describe("UpdateUser", func() {
		It("should update a user", func() {
			body, err := json.Marshal(request)
			Expect(err).To(BeNil())

			ginCtx.Request, err = http.NewRequest("PUT", "/v1/users/:user_id", bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			var response = &models.UserResponse{
				ID:        1,
				Firstname: "tmostowasherenot",
				Lastname:  "tmostowasherenot",
				Email:     "tmostowashere@tmostowashere.com",
				BookList:  nil,
				Role:      "superuser",
			}
			fakeUserer.UpdateUserReturns(response, nil)

			router.ServeHTTP(w, ginCtx.Request)

			actualBody, _ := io.ReadAll(w.Result().Body)
			Expect(w.Code).To(Equal(http.StatusCreated), "Expected HTTP status to be OK, but got %d", w.Code)
			Expect(w.Body.String()).To(Equal(string(actualBody)), "Unexpected response body: %s", w.Body.String())
		})
	})

	Describe("DeleteUser", func() {
		It("should delete a user successfully", func() {
			var err error

			ginCtx.Request, err = http.NewRequest("DELETE", "/v1/users/:user_id/2", nil)
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeUserer.DeleteUserReturns(2, nil)

			router.ServeHTTP(w, ginCtx.Request)

			Expect(w.Code).To(Equal(http.StatusOK), "Expected HTTP status to be OK, but got %d", w.Code)
			Expect(w.Body.String()).To(ContainSubstring("User deleted successfully"))
		})
	})

	Describe("GetAllUsers", func() {
		It("should retrieve all users", func() {
			var err error

			ginCtx.Request, err = http.NewRequest("GET", "/v1/users", nil)
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeUserer.GetAllUsersReturns([]models.UserResponse{
				{
					ID:        1,
					Firstname: "User1",
					Lastname:  "User1",
					Email:     "user1@user.com",
					BookList:  nil,
					Role:      "user",
				},
				{
					ID:        2,
					Firstname: "User2",
					Lastname:  "User2",
					Email:     "user2@user.com",
					BookList:  nil,
					Role:      "user",
				},
			}, nil)

			router.ServeHTTP(w, ginCtx.Request)

			actualBody, _ := io.ReadAll(w.Result().Body)
			Expect(w.Code).To(Equal(http.StatusOK), "Expected HTTP status OK")
			Expect(w.Body.String()).To(Equal(string(actualBody)), "Unexpected response body: %s", w.Body.String())
		})
	})

	Describe("GetUser", func() {
		It("should retrieve a user", func() {
			var err error

			response := &models.UserResponse{
				ID:        1,
				Firstname: "User1",
				Lastname:  "User1",
				Email:     "user1@user.com",
				BookList:  nil,
				Role:      "user",
			}

			ginCtx.Request, err = http.NewRequest("GET", "/v1/users/1", nil)
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeUserer.GetUserReturns(response, nil)

			router.ServeHTTP(w, ginCtx.Request)

			Expect(w.Code).To(Equal(http.StatusOK), "Expected HTTP status OK")
			actualBody, _ := io.ReadAll(w.Result().Body)
			Expect(w.Body.String()).To(Equal(string(actualBody)), "Unexpected response body: %s", w.Body.String())
		})
	})
})
