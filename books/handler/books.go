package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"library/books/models"
	"library/books/repository"
	"library/pkg"
	"library/pkg/postgres"
	"library/pkg/utils"
	"net/http"
)

type BookerHandler interface {
	AddBook(c *gin.Context)
	UpdateBook(c *gin.Context)
	GetBook(c *gin.Context)
	GetAllBooks(c *gin.Context)
	DeleteBook(c *gin.Context)
}

type BookHandler struct {
	ctx            context.Context
	bookRepository repository.BookerRepository
}

func NewBookHandler(ctx context.Context, booker repository.BookerRepository) BookerHandler {
	return &BookHandler{
		ctx:            ctx,
		bookRepository: booker,
	}
}

func (b *BookHandler) AddBook(c *gin.Context) {
	var book models.BookRequest
	var err error

	log := utils.GetLogger(b.ctx)

	book.UserID.ID = c.GetInt("userID")

	if err = c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isExisting, err := repository.ValidateISBNExists(log, book.ISBN, b.bookRepository.GetDBPool())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isExisting {
		errorMessage := "Book with this ISBN already exists"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = pkg.CheckDate(log, &book)
	if err != nil {
		errorMessage := "Cannot add a book with a future publication date"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.AddBook(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Book added successfully": book})
}

func (b *BookHandler) UpdateBook(c *gin.Context) {
	var book models.BookRequest
	var err error

	log := utils.GetLogger(b.ctx)

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = c.GetInt("book_id")

	exists, err := postgres.CheckIDExists("books", book.ID, b.bookRepository.GetDBPool())
	if err != nil {
		errorMessage := fmt.Sprintf("Checking book ID error: %v", book.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	if !exists {
		errorMessage := fmt.Sprintf("Book ID doesn't exists: %v", book.ID)
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	book.UserID.ID = c.GetInt("userID")

	isAssigned, err := repository.IsAssigned(log, book.ID, book.UserID.ID, b.bookRepository.GetDBPool())
	if err != nil {
		errorMessage := "Error checking book assignment"
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	if !isAssigned {
		errorMessage := "User is not the owner of the book"
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
		return
	}

	err = pkg.CheckDate(log, &book)
	if err != nil {
		errorMessage := "Cannot add a book with a future publication date"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.UpdateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (b *BookHandler) GetBook(c *gin.Context) {
	var book models.BookResponse
	var err error

	book.ID = c.GetInt("book_id")

	err = b.bookRepository.GetBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (b *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := b.bookRepository.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (b *BookHandler) DeleteBook(c *gin.Context) {
	log := utils.GetLogger(b.ctx)

	bookID := c.GetInt("book_id")

	exists, err := postgres.CheckIDExists("user_book", bookID, b.bookRepository.GetDBPool())
	if err != nil {
		errorMessage := fmt.Sprintf("Checking book ID error: %v", bookID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	if !exists {
		errorMessage := fmt.Sprintf("Book ID doesn't exists: %v", bookID)
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	userID := c.GetInt("userID")

	isAssigned, err := repository.IsAssigned(log, bookID, userID, b.bookRepository.GetDBPool())
	if err != nil {
		errorMessage := "Error checking book assignment"
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	if !isAssigned {
		errorMessage := "User is not the owner of the book"
		c.JSON(http.StatusUnauthorized, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.DeleteBook(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
