package domain

import (
	"time"
)

type Book struct {
	ID          int        `json:"-"`                         //swagger:ignore
	UUID        string     `json:"uuid" swaggerignore:"true"` // Ignora no input
	Title       string     `json:"title" example:"1984"`      // @example 1984
	Authors     []*Author  `json:"authors"`
	Description string     `json:"description" example:"Livro conta a história...."` //@example Livro conta a história....
	CreatedAt   *time.Time `json:"-,omitempty"`                                      //swagger:ignore
} //@name Book
type BookCreateRequest struct {
	Title       string `json:"title" binding:"required,min=2,max=100" example:"1984"`
	Description string `json:"description,omitempty" example:"A dystopian novel" binding:"max=500"`
	AuthorIDs   []int  `json:"authorIds" example:"1,2,3"`
} //@name AuthorRequest

type BookCreateResponse struct {
	UUID        string    `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string    `json:"title" example:"1984"`
	Description string    `json:"description,omitempty" example:"A dystopian novel"`
	Authors     []Author  `json:"authors"`
	CreatedAt   time.Time `json:"createdAt" example:"2023-01-01T00:00:00Z"`
} // @name BookCreateResponse
type BookFilters struct {
	Title         string
	AuthorName    string // Only used in WithAuthors version
	Sort          string // "title", "created_at"
	SortDirection string // "ASC", "DESC"
	Limit         int    // 10, 25, 50...
	Offset        int    // (page-1)*limit
} //@nome BookFilters
type BookListResponse struct {
	Data  []Book `json:"data"`
	Total int    `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
} //@name BookListResponse

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error   string `json:"error"`             // Error message
	Details string `json:"details,omitempty"` // Additional details (debug only)
} // @name ErrorResponse
