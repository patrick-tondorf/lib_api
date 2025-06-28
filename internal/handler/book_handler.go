package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patrick-tondorf/lib_api/internal/domain"
	"github.com/patrick-tondorf/lib_api/internal/repository"
)

// BookHandler defines the book handler methods
type BookHandler struct {
	Repo *repository.BookRepository
}

// NewBookHandler creates a new BookHandler.
func NewBookHandler(repo *repository.BookRepository) *BookHandler {
	return &BookHandler{Repo: repo}
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the library
// @Tags books
// @Accept  json
// @Produce  json
// @Param   book  body  domain.Book  true  "Book data"
// @Success 201 {object} domain.Book
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book domain.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Repo.CreateBook(c.Request.Context(), &book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// GetBooks godoc
// @Summary List all books
// @Description Get all books from the library
// @Tags books
// @Produce  json
// @Success 200 {array} domain.Book
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.Repo.GetBooks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update book details by ID
// @Tags books
// @Accept  json
// @Produce  json
// @Param   id    path  int         true  "Book ID"
// @Param   book  body  domain.Book true  "Updated book data"
// @Success 200 {object} domain.Book
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "Book not found"
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	// Implementação...
}
