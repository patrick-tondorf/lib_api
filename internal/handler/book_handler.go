package handler

import (
	"net/http"
	"slices"

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
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param   book  body  domain.BookCreateRequest  true  "Book data"
// @Example request({"title":"1984","description":"A dystopian novel","authors":[{"id":1},{"id":2}]})
// @Success 201 {object} domain.Book
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book domain.BookCreateRequest
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(), // Mostra detalhes do erro de validação
		})
		return
	}
	/*	var book domain.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}*/

	if err := h.Repo.CreateBook(c.Request.Context(), book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// GetBooks godoc
// @Summary List books with pagination and filters
// @Description Get paginated list of books with optional filters. Choose between basic version or with authors.
// @Tags books
// @Security BearerAuth
// @Produce json
// @Param title        query string  false "Filter by book title (partial match, case insensitive)"
// @Param author       query string  false "Filter by author name (only works when with_authors=true)"
// @Param with_authors query boolean false "Include full author information in response"
// @Param sort         query string  false "Sort field" Enums(title, created_at) default(title)
// @Param sort_dir     query string  false "Sort direction" Enums(ASC, DESC) default(ASC)
// @Param page         query int     false "Page number" default(1) minimum(1)
// @Param limit        query int     false "Items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} BookListResponse
// @Failure 400 {object} ErrorResponse "Invalid parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	// Parse common filters
	filters := domain.BookFilters{
		Title:         c.Query("title"),
		AuthorName:    c.Query("author"),
		Sort:          c.DefaultQuery("sort", "title"),
		SortDirection: c.DefaultQuery("sort_dir", "ASC"),
		Limit:         clamp(c.GetInt("limit"), 1, 100),
		Offset:        (clamp(c.GetInt("page"), 1, 1000) - 1) * clamp(c.GetInt("limit"), 1, 100),
	}

	// Validate sort
	if !slices.Contains([]string{"title", "created_at"}, filters.Sort) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort field"})
		return
	}

	// Choose repository method based on query param
	var (
		books []domain.Book
		total int
		err   error
	)

	if c.Query("with_authors") == "true" {
		books, total, err = h.Repo.GetBooksWithAuthors(c.Request.Context(), filters)
	} else {
		books, total, err = h.Repo.GetBooksBasic(c.Request.Context(), filters)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.BookListResponse{
		Data:  books,
		Total: total,
		Page:  filters.Offset/filters.Limit + 1,
		Limit: filters.Limit,
	})
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update book details by ID
// @Tags books
// @Security BearerAuth
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

// GetBook godoc
// @Summary Get a book by ID
// @Description Retrieve a book by its ID
// @Tags books
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} domain.Book
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	/*	id := c.Param("id")

		book, err := h.Repo.GetBookByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		c.JSON(http.StatusOK, book)*/
}

// Helper function to clamp values
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
