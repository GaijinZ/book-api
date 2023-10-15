package handler

import (
	"library/books/models"
	"library/books/repository"
	"library/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	middleware "library/pkg/middleware"
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

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := middleware.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	book.UserID.ID = userID

	if err = c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isExisting, err := repository.ValidateISBNExists(book.ISBN, b.bookRepository.DBPool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isExisting {
		errorMessage := "Book with this ISBN already exists"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = pkg.CheckDate(&book)
	if err != nil {
		errorMessage := "Cannot add a book with a future publication date"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	err = b.bookRepository.AddBook(userID, book)
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
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := pkg.CheckIDExists("books", bookID, b.bookRepository.DBPool)
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

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := middleware.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isAssigned, err := repository.IsAssigned(bookID, userID, b.bookRepository.DBPool)
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

	err = pkg.CheckDate(&book)
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
	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		errorMessage := "Wrong user ID: " + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}

	exists, err := pkg.CheckIDExists("books", bookID, b.bookRepository.DBPool)
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

	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := middleware.VerifyJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := strconv.Atoi(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isAssigned, err := repository.IsAssigned(bookID, userID, b.bookRepository.DBPool)
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
