package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"library/books/handler"
	"library/books/models"
	"library/books/repository/repositoryfakes"
	"library/books/server"
	"library/pkg/logger"
	"library/pkg/middleware"
	userModel "library/users/models"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Book API Test", func() {
	var (
		fakeBooker      *repositoryfakes.FakeBookerRepository
		fakeBookHandler handler.BookerHandler
		w               *httptest.ResponseRecorder
		ginCtx          *gin.Context
		router          *gin.Engine
		user            *userModel.User
		enricher        func(req *http.Request) *http.Request
		cookie          *http.Cookie
	)

	JustBeforeEach(func() {
		w = httptest.NewRecorder()
		ginCtx, _ = gin.CreateTestContext(w)

		ctx := context.Background()
		ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))

		fakeBooker = &repositoryfakes.FakeBookerRepository{}
		fakeBookHandler = handler.NewBookHandler(ctx, fakeBooker)

		router = server.NewRouter(fakeBookHandler)

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

	Describe("AddBook", func() {
		It("should add a book", func() {
			var request = &models.BookRequest{
				ID:            1,
				Name:          "tmostowashere",
				DatePublished: "2022-01-01",
				ISBN:          "12345679",
				PageCount:     123,
				Author: models.AuthorRequest{
					ID:   1,
					Name: "tmostowashere",
				},
			}

			body, err := json.Marshal(request)
			Expect(err).To(BeNil())

			ginCtx.Request, err = http.NewRequest("POST", "/v1/books/:user_id/books", bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeBooker.AddBookReturns(1, nil)

			router.ServeHTTP(w, ginCtx.Request)

			Expect(w.Code).To(Equal(http.StatusCreated), "Expected HTTP status Created")
			args := fakeBooker.AddBookArgsForCall(0)
			Expect(args.ID).To(Equal(1), "Expected AddBOok to be called with the correct arguments")
		})
	})

	Describe("UpdateBook", func() {
		It("should update a book", func() {
			var request = &models.BookRequest{
				ID:            1,
				Name:          "tmostowashere",
				DatePublished: "2022-01-01",
				ISBN:          "12345679",
				PageCount:     123,
				Author: models.AuthorRequest{
					ID:   1,
					Name: "tmostowashere",
				},
			}

			body, err := json.Marshal(request)
			Expect(err).To(BeNil())

			ginCtx.Request, err = http.NewRequest("PUT", "/v1/books/:user_id/books/1", bytes.NewBuffer(body))
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			var response = &models.BookResponse{
				ID:            1,
				Name:          "tmostowasherenot",
				DatePublished: "2022-01-01",
				ISBN:          "12345679",
				PageCount:     123,
				Author: models.AuthorResponse{
					ID:   1,
					Name: "tmostowashere",
				},
			}

			fakeBooker.UpdateBookReturns(response, nil)

			router.ServeHTTP(w, ginCtx.Request)

			Expect(w.Code).To(Equal(http.StatusCreated), "Expected HTTP status Created")
			actualBody, _ := io.ReadAll(w.Result().Body)
			Expect(w.Body.String()).To(Equal(string(actualBody)), "Unexpected response body: %s", w.Body.String())
		})
	})

	Describe("DeleteBoook", func() {
		It("should delete a book", func() {
			var err error

			ginCtx.Request, err = http.NewRequest("DELETE", "/v1/books/:user_id/books/2", nil)
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeBooker.DeleteBookReturns(2, nil)

			router.ServeHTTP(w, ginCtx.Request)

			Expect(w.Code).To(Equal(http.StatusOK), "Expected HTTP status to be OK, but got %d", w.Code)
			Expect(w.Body.String()).To(ContainSubstring("Book deleted successfully"))
		})
	})

	Describe("GetAllBooks", func() {
		It("should retrieve all books", func() {
			var err error

			ginCtx.Request, err = http.NewRequest("GET", "/v1/books/:user_id/books", nil)
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeBooker.GetAllBooksReturns([]models.BookResponse{
				{
					ID:            1,
					Name:          "Book1",
					DatePublished: "Book1",
					ISBN:          "1234679",
					PageCount:     123,
					UserID: userModel.User{
						ID: 1,
					},
					Author: models.AuthorResponse{
						ID:   1,
						Name: "Author1",
					},
				},
				{
					ID:            2,
					Name:          "Book2",
					DatePublished: "Book2",
					ISBN:          "987654321",
					PageCount:     321,
					UserID: userModel.User{
						ID: 2,
					},
					Author: models.AuthorResponse{
						ID:   2,
						Name: "Author2",
					},
				},
			}, nil)

			router.ServeHTTP(w, ginCtx.Request)

			actualBody, _ := io.ReadAll(w.Result().Body)
			Expect(w.Code).To(Equal(http.StatusOK), "Expected HTTP status OK")
			Expect(w.Body.String()).To(Equal(string(actualBody)), "Unexpected response body: %s", w.Body.String())
		})
	})

	Describe("GetBook", func() {
		It("should retrieve a book", func() {
			var err error

			response := &models.BookResponse{
				ID:            1,
				Name:          "Book1",
				DatePublished: "Book1",
				ISBN:          "1234679",
				PageCount:     123,
				UserID: userModel.User{
					ID: 1,
				},
				Author: models.AuthorResponse{
					ID:   1,
					Name: "Author1",
				},
			}

			ginCtx.Request, err = http.NewRequest("GET", "/v1/books/:user_id/books/1", nil)
			Expect(err).To(BeNil())
			ginCtx.Request.AddCookie(cookie)
			ginCtx.Request = enricher(ginCtx.Request)

			fakeBooker.GetBookReturns(response, nil)

			router.ServeHTTP(w, ginCtx.Request)

			actualBody, _ := io.ReadAll(w.Result().Body)
			Expect(w.Code).To(Equal(http.StatusOK), "Expected HTTP status OK")
			Expect(w.Body.String()).To(Equal(string(actualBody)), "Unexpected response body: %s", w.Body.String())
		})
	})
})
