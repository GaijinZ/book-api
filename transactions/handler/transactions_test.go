package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"library/pkg/logger"
	"library/pkg/middleware"
	"library/transactions/handler"
	"library/transactions/models"
	"library/transactions/repository/repositoryfakes"
	"library/transactions/server"
	userModel "library/users/models"

	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Transaction API Test", func() {
	var (
		fakeTransactioner      *repositoryfakes.FakeTransactionerRepository
		fakeTransactionHandler handler.TransactionerHandler
		w                      *httptest.ResponseRecorder
		ginCtx                 *gin.Context
		router                 *gin.Engine
		enricher               func(req *http.Request) *http.Request
		cookie                 *http.Cookie
		user                   *userModel.User
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		ginCtx, _ = gin.CreateTestContext(w)

		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))

		fakeTransactioner = &repositoryfakes.FakeTransactionerRepository{}
		fakeTransactionHandler = handler.NewTransactionHandler(ctx, fakeTransactioner)

		router = server.NewRouter(fakeTransactionHandler)

		user = &userModel.User{
			ID:    1,
			Email: "tmostowashere@tmostowashere.com",
			Role:  "superuser",
		}

		token, _ := middleware.GenerateJWT(*user)
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

	Describe("Buy Book", func() {
		It("should add a book to user", func() {
			var request = &models.TransactionResponse{
				Quantity: 1,
			}

			body, err := json.Marshal(request)
			Expect(err).To(BeNil())

			ginCtx.Request, err = http.NewRequest("POST", "/v1/transactions/buy-book/1", bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeTransactioner.BuyBookReturns(1, nil)

			router.ServeHTTP(w, ginCtx.Request)

			Expect(w.Code).To(Equal(http.StatusCreated), "Expected HTTP status Created")
			actualBody, _ := io.ReadAll(w.Result().Body)
			Expect(w.Body.String()).To(Equal(string(actualBody)), "Unexpected response body: %s", w.Body.String())
		})
	})

	Describe("Transaction History", func() {
		It("should return user transaction history", func() {
			var err error

			ginCtx.Request, err = http.NewRequest("POST", "/v1/transactions/transactions", nil)
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			transactionResponse := []models.UserTransactionResponse{
				{
					BookList:        nil,
					UserID:          1,
					Quantity:        2,
					TransactionDate: time.Time{},
				},
			}

			fakeTransactioner.TransactionHistoryReturns(transactionResponse, nil)

			router.ServeHTTP(w, ginCtx.Request)

			Expect(w.Code).To(Equal(http.StatusOK), "Expected HTTP status OK")
			actualBody, _ := io.ReadAll(w.Result().Body)
			Expect(w.Body.String()).To(Equal(string(actualBody)), "Unexpected response body: %s", w.Body.String())
		})
	})
})
