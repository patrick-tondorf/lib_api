package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/patrick-tondorf/lib_api/internal/domain"

	"github.com/jackc/pgx/v5"
)

type BookRepository struct {
	DB *pgx.Conn
}

func NewBookRepository(db *pgx.Conn) *BookRepository {
	return &BookRepository{DB: db}
}

// Create
func (r *BookRepository) CreateBook(ctx context.Context, book *domain.Book) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO books (title,author) VALUES ($1,$2)", book.Title, book.Author)
	return err
}

// Get ALl
func (r *BookRepository) GetBooks(ctx context.Context) ([]domain.Book, error) {
	log.Println("Attempting to query books from database")
	rows, err := r.DB.Query(ctx, "SELECT id, title, author, created_at FROM books")
	if err != nil {
		log.Printf("Database query error: %v\n", err) // Adicione esta linha
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var b domain.Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.CreatedAt)
		if err != nil {
			log.Printf("Row scan error: %v\n", err) // Adicione esta linha
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v\n", err) // Adicione esta linha
		return nil, fmt.Errorf("rows error: %w", err)
	}

	log.Printf("Successfully retrieved %d books\n", len(books))
	return books, nil
}
