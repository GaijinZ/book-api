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

// AddBook adds a new book to the database.
//
//	@Summary		Add a new book
//	@Description	Adds a new book to the database.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string										true	"JWT Token"
//	@Param			book			body		models.BookRequest									true	"Book object to be added"
//	@Success		201
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/v1/books [post]
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

// UpdateBook updates an existing book in the database.
//
//	@Summary		Update an existing book
//	@Description	Updates an existing book in the database.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"JWT Token"
//	@Param			book			body		models.BookRequest												true	"Updated book object"
//	@Success		201
//	@Failure		400
//	@Failure		401
//	@Failure		500
//	@Router			/v1/books [put]
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

// GetBook retrieves a book by its ID from the database.
//
//	@Summary		Retrieve a book by ID
//	@Description	Retrieves a book by its ID from the database.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"JWT Token"
//	@Param			book_id			path		int									true	"Book ID"
//	@Success		200
//	@Failure		500
//	@Router			/v1/books/{book_id} [get]
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

// GetAllBooks retrieves all books from the database.
//
//	@Summary		Retrieve all books
//	@Description	Retrieves all books from the database.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"JWT Token"
//	@Success		200
//	@Failure		400
//	@Router			/v1/books [get]
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

// DeleteBook deletes a book from the database.
//
//	@Summary		Delete a book
//	@Description	Deletes a book from the database.
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			bookID			path		int							true	"Book ID to delete"
//	@Param			Authorization	header		string						true	"JWT Token"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/v1/books/{bookID} [delete]
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
