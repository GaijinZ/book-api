package handler

import (
	"context"
	"library/books/models"
	"library/books/repository"
	"library/pkg"
	"library/pkg/logger"
	"library/pkg/postgres"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	var book models.Book
	var err error
	log := b.ctx.Value("logger").(logger.Logger)

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

	err = b.bookRepository.AddBook(book.UserID.ID, book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Book added successfully": book})
}

func (b *BookHandler) UpdateBook(c *gin.Context) {
	var book models.Book
	log := b.ctx.Value("logger").(logger.Logger)

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		errorMessage := "Wrong book ID: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := postgres.CheckIDExists("books", bookID, b.bookRepository.GetDBPool())
	if err != nil {
		errorMessage := "Checking book ID error: " + string(rune(bookID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	if !exists {
		errorMessage := "Book ID doesn't exists: " + string(rune(bookID))
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	book.UserID.ID = c.GetInt("userID")

	isAssigned, err := repository.IsAssigned(log, bookID, book.UserID.ID, b.bookRepository.GetDBPool())
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

	err = b.bookRepository.UpdateBook(bookID, book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (b *BookHandler) GetBook(c *gin.Context) {
	var book models.Book

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		errorMessage := "Wrong book ID: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.GetBook(bookID, &book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (b *BookHandler) GetAllBooks(c *gin.Context) {
	users, err := b.bookRepository.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"books": users})
}

func (b *BookHandler) DeleteBook(c *gin.Context) {
	log := b.ctx.Value("logger").(logger.Logger)

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		errorMessage := "Wrong user ID: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	exists, err := postgres.CheckIDExists("books", bookID, b.bookRepository.GetDBPool())
	if err != nil {
		errorMessage := "Checking book ID error: " + string(rune(bookID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	if !exists {
		errorMessage := "Book ID doesn't exists: " + string(rune(bookID))
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
