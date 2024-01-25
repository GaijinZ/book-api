package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"library/books/models"
	"library/books/repository"
	"library/pkg"
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
		log.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = pkg.CheckPublishedDate(log, &book)
	if err != nil {
		log.Errorf("checking published date error: %v", err)
		errorMessage := "Cannot add a book with a future publication date"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	bookID, err := b.bookRepository.AddBook(&book)
	if err != nil {
		log.Errorf("Add Book repository error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Book added successfully: %v", &book)
	c.JSON(http.StatusCreated, gin.H{"Book added successfully": bookID})
}

func (b *BookHandler) UpdateBook(c *gin.Context) {
	var book models.BookRequest
	var bookResponse *models.BookResponse

	log := utils.GetLogger(b.ctx)

	if err := c.ShouldBindJSON(&book); err != nil {
		log.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := pkg.CheckPublishedDate(log, &book)
	if err != nil {
		log.Errorf("checking published date error: %v", err)
		errorMessage := "Cannot add a book with a future publication date"
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	book.ID = c.GetInt("bookID")
	book.UserID.ID = c.GetInt("userID")

	bookResponse, err = b.bookRepository.UpdateBook(&book)
	if err != nil {
		log.Errorf("Update book repository error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Book updated successfully: %v", &book)
	c.JSON(http.StatusCreated, gin.H{"Book updated successfully": bookResponse})
}

func (b *BookHandler) GetBook(c *gin.Context) {
	var book *models.BookResponse
	var err error

	log := utils.GetLogger(b.ctx)

	bookID := c.GetInt("bookID")

	book, err = b.bookRepository.GetBook(bookID)
	if err != nil {
		log.Errorf("Get Book repository error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Book: %v", book)
	c.JSON(http.StatusOK, gin.H{"book": book})
}

func (b *BookHandler) GetAllBooks(c *gin.Context) {
	log := utils.GetLogger(b.ctx)

	books, err := b.bookRepository.GetAllBooks()
	if err != nil {
		log.Errorf("Get All Books repository error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Books: %v", books)
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (b *BookHandler) DeleteBook(c *gin.Context) {
	log := utils.GetLogger(b.ctx)

	bookID := c.GetInt("bookID")
	userID := c.GetInt("userID")

	deletedID, err := b.bookRepository.DeleteBook(bookID, userID)
	if err != nil {
		log.Errorf("Delete Book repository error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Infof("Book deleted successfully, id: %v", deletedID)
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
