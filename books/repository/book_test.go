package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"library/books/models"
	"library/books/repository"
	"library/pkg/logger"
	"library/pkg/postgres"
	userModel "library/users/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API Test", func() {
	var (
		bookRepo    repository.BookerRepository
		bookRequest *models.BookRequest
		//bookResponse *models.BookResponse
		mock   sqlmock.Sqlmock
		fakeDB *postgres.DB
		//err          error
	)

	JustBeforeEach(func() {
		var ctx context.Context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctx = context.WithValue(ctx, "logger", logger.NewLogger(2))

		fakeDB, _ = postgres.NewFakeDB(ctx)

		bookRepo = repository.NewBookRepository(ctx, *fakeDB)

		bookRequest = &models.BookRequest{
			ID:            1,
			Name:          "tmosto",
			DatePublished: "tmosto",
			ISBN:          "1234679",
			PageCount:     123,
			UserID: userModel.User{
				ID: 1,
			},
			Author: models.AuthorRequest{
				ID:   1,
				Name: "tmosto",
			},
		}

		mock = fakeDB.GetMock()
	})

	AfterEach(func() {
		fakeDB.Close()
	})

	Describe("AddUser", func() {
		It("should add a user successfully", func() {
			mock.ExpectQuery(repository.CheckISBN).
				WithArgs(bookRequest.ISBN).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectQuery(repository.SelectAuthor).
				WithArgs(bookRequest.Author.Name).
				WillReturnError(sql.ErrNoRows)

			mock.ExpectExec(repository.InsertAuthor).
				WithArgs(bookRequest.Author.Name).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectExec(repository.InsertBook).
				WithArgs(bookRequest.Name, bookRequest.DatePublished, bookRequest.ISBN, bookRequest.PageCount, bookRequest.UserID.ID, 1).
				WillReturnResult(sqlmock.NewResult(1, 1))

			bookID, err := bookRepo.AddBook(bookRequest)
			Expect(err).To(BeNil())
			Expect(bookID).To(Equal(1), "Expected AddUser to be called with the correct arguments")
		})

		It("should return an error if ISBN already exists", func() {
			mock.ExpectQuery(repository.CheckISBN).
				WithArgs(bookRequest.ISBN).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

			_, err := bookRepo.AddBook(bookRequest)

			Expect(err).NotTo(BeNil())
			Expect(err).To(MatchError("book with this ISBN already exists"))
		})

		It("should return an error if author retrieval or creation fails", func() {
			mock.ExpectQuery(repository.CheckISBN).
				WithArgs(bookRequest.ISBN).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectExec(repository.InsertAuthor).
				WithArgs(bookRequest.Author.Name).
				WillReturnError(errors.New("author error"))

			_, err := bookRepo.AddBook(bookRequest)

			Expect(err).NotTo(BeNil())
			Expect(err).To(MatchError("author error"))
		})

		It("should return an error if book insertion fails", func() {
			mock.ExpectQuery(repository.CheckISBN).
				WithArgs(bookRequest.ISBN).
				WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

			mock.ExpectQuery(repository.SelectAuthor).
				WithArgs(bookRequest.Author.Name).
				WillReturnError(sql.ErrNoRows)

			mock.ExpectExec(repository.InsertBook).
				WithArgs(bookRequest)
		})

		Describe("UpdateBook", func() {
			Context("when book exists and user is assigned", func() {
				It("should update the book", func() {
					book := &models.BookRequest{
						ID:            1,
						Name:          "Updated Book",
						DatePublished: "2022-01-01",
						ISBN:          "1234567890",
						PageCount:     200,
						Author: models.AuthorRequest{
							Name: "Author Name",
						},
						UserID: userModel.User{ID: 1},
					}

					query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM books WHERE id=$1)")
					mock.ExpectQuery(query).
						WithArgs(book.ID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

					mock.ExpectQuery(repository.IsAssigned).
						WithArgs(book.ID, book.UserID.ID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

					mock.ExpectQuery(repository.SelectAuthor).
						WithArgs(book.Author.Name).
						WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

					mock.ExpectExec(repository.UpdateBook).
						WithArgs(book.Name, book.DatePublished, book.ISBN, book.PageCount, 1, book.ID).
						WillReturnResult(sqlmock.NewResult(1, 1))

					bookResponse, err := bookRepo.UpdateBook(book)
					Expect(err).NotTo(HaveOccurred())
					Expect(bookResponse).NotTo(BeNil())
				})
			})

			Context("when book does not exist", func() {
				It("should return an error", func() {
					book := &models.BookRequest{
						ID:            2,
						Name:          "Updated Book",
						DatePublished: "2022-01-01",
						ISBN:          "1234567890",
						PageCount:     200,
						Author: models.AuthorRequest{
							Name: "Author Name",
						},
						UserID: userModel.User{ID: 1},
					}

					query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM books WHERE id=$1)")
					mock.ExpectQuery(query).
						WithArgs(book.ID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

					_, err := bookRepo.UpdateBook(book)
					Expect(err).To(HaveOccurred())
					expectedErrorMessage := fmt.Sprintf("Book ID doesn't exists: %d", book.ID)
					Expect(err.Error()).To(Equal(expectedErrorMessage))
				})
			})

			Context("when book exists but user is not assigned", func() {
				It("should return an error", func() {
					book := &models.BookRequest{
						ID:            1,
						Name:          "Updated Book",
						DatePublished: "2022-01-01",
						ISBN:          "1234567890",
						PageCount:     200,
						Author: models.AuthorRequest{
							Name: "Author Name",
						},
						UserID: userModel.User{ID: 2},
					}

					query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM books WHERE id=$1)")
					mock.ExpectQuery(query).
						WithArgs(book.ID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

					mock.ExpectQuery(repository.IsAssigned).
						WithArgs(book.ID, book.UserID.ID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

					_, err := bookRepo.UpdateBook(book)
					Expect(err).To(HaveOccurred())
					expectedErrorMessage := "user is not the owner of the book"
					Expect(err.Error()).To(Equal(expectedErrorMessage))
				})
			})
		})

		Describe("GetBook", func() {
			Context("when the book exists", func() {
				It("should return the book details", func() {
					expectedBook := &models.BookResponse{
						Name:          "Test Book",
						DatePublished: "2022-01-01",
						ISBN:          "1234567890",
						PageCount:     200,
						Author: models.AuthorResponse{
							Name: "Author Name",
						},
					}
					bookID := 123

					mock.ExpectQuery(repository.GetBook).
						WithArgs(bookID).
						WillReturnRows(sqlmock.NewRows([]string{"name", "date_published", "isbn", "page_count", "author_name"}).
							AddRow(expectedBook.Name, expectedBook.DatePublished, expectedBook.ISBN, expectedBook.PageCount, expectedBook.Author.Name))

					bookResponse, err := bookRepo.GetBook(bookID)

					Expect(err).ToNot(HaveOccurred())
					Expect(bookResponse).To(Equal(expectedBook))
				})
			})

			Context("when the book does not exist", func() {
				It("should return an error", func() {
					expectedBook := &models.BookResponse{}
					bookID := 456

					mock.ExpectQuery(repository.GetBook).
						WithArgs(bookID).
						WillReturnError(sql.ErrNoRows)

					bookResponse, err := bookRepo.GetBook(bookID)

					Expect(err).To(HaveOccurred())
					Expect(bookResponse).To(Equal(expectedBook))
				})
			})
		})

		Describe("GetAllBooks", func() {
			Context("when books exist", func() {
				It("should return a list of books", func() {
					expectedBooks := []models.BookResponse{
						{Name: "Book1", DatePublished: "2022-01-01", ISBN: "1234567890", PageCount: 200, Author: models.AuthorResponse{Name: "Author1"}},
						{Name: "Book2", DatePublished: "2022-02-01", ISBN: "0987654321", PageCount: 250, Author: models.AuthorResponse{Name: "Author2"}},
					}

					rows := sqlmock.NewRows([]string{"name", "date_published", "isbn", "page_count", "author_name"})
					for _, book := range expectedBooks {
						rows.AddRow(book.Name, book.DatePublished, book.ISBN, book.PageCount, book.Author.Name)
					}

					mock.ExpectQuery(repository.GetAllBooks).
						WillReturnRows(rows)

					books, err := bookRepo.GetAllBooks()

					Expect(err).NotTo(HaveOccurred())

					Expect(books).To(Equal(expectedBooks))
				})
			})

			Context("when no books exist", func() {
				It("should return an empty list", func() {
					mock.ExpectQuery(repository.GetAllBooks).
						WillReturnRows(sqlmock.NewRows([]string{}))

					books, err := bookRepo.GetAllBooks()

					Expect(err).NotTo(HaveOccurred())

					Expect(books).To(BeEmpty())
				})
			})
		})

		Describe("DeleteBook", func() {
			Context("when the book exists and the user is the owner", func() {
				It("should delete the book", func() {
					bookID := 123
					userID := 456

					query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM user_book WHERE id=$1)")
					mock.ExpectQuery(query).
						WithArgs(bookID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

					mock.ExpectQuery(repository.IsAssigned).
						WithArgs(bookID, userID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

					mock.ExpectExec(repository.DeleteBook).
						WithArgs(bookID).
						WillReturnResult(sqlmock.NewResult(0, 1))

					deletedBookID, err := bookRepo.DeleteBook(bookID, userID)

					Expect(err).NotTo(HaveOccurred())

					Expect(deletedBookID).To(Equal(bookID))
				})
			})

			Context("when the book does not exist", func() {
				It("should return an error", func() {
					bookID := 123
					userID := 456

					query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM user_book WHERE id=$1)")
					mock.ExpectQuery(query).
						WithArgs(bookID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

					_, err := bookRepo.DeleteBook(bookID, userID)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("Book ID doesn't exists: %v", bookID)))
				})
			})

			Context("when the user is not the owner of the book", func() {
				It("should return an error", func() {
					bookID := 123
					userID := 456

					query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM user_book WHERE id=$1)")
					mock.ExpectQuery(query).
						WithArgs(bookID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

					mock.ExpectQuery(repository.IsAssigned).
						WithArgs(bookID, userID).
						WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

					_, err := bookRepo.DeleteBook(bookID, userID)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("user is not the owner of the book"))
				})
			})
		})
	})
})
