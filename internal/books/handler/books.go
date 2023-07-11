package handler

import (
	"library/internal/books/models"
	"library/internal/books/repository"
	"library/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Booker interface {
	AddBook(c *gin.Context)
	UpdateBook(c *gin.Context)
	GetBook(c *gin.Context)
	GetAllBooks(c *gin.Context)
	DeleteBook(c *gin.Context)
}

type BookHandler struct {
	bookRepository *repository.BookRepository
}

func NewBookHandler(bookRepository *repository.BookRepository) *BookHandler {
	return &BookHandler{
		bookRepository: bookRepository,
	}
}

func (b *BookHandler) AddBook(c *gin.Context) {
	var book models.Book
	var err error

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		errorMessage := "Invalid user ID: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	book.UserID.ID = userID

	if err = c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error11": err.Error()})
		return
	}

	err = utils.CheckDate(&book)
	if err != nil {
		errorMessage := "Cannot add a book with a future publication date"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.AddBook(userID, book, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Book added successfully": book})
}

func (b *BookHandler) UpdateBook(c *gin.Context) {
	var book models.Book

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		errorMessage := "Wrong book ID: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = utils.CheckDate(&book)
	if err != nil {
		errorMessage := "Cannot add a book with a future publication date"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.UpdateBook(bookID, book, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Book updated successfully": book})
}

func (b *BookHandler) GetBook(c *gin.Context) {
	var book models.Book

	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		errorMessage := "Wrong book ID: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.GetBook(bookID, &book, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (b *BookHandler) GetAllBooks(c *gin.Context) {
	users, err := b.bookRepository.GetAllBooks(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (b *BookHandler) DeleteBook(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		errorMessage := "Wrong user ID: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.DeleteBook(bookID, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
