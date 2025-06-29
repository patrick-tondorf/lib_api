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
func (r *BookRepository) CreateBook(ctx context.Context, req domain.BookCreateRequest) error {
	// Verificar se todos os autores existem antes de começar a transação
	if len(req.AuthorIDs) == 0 {
		return fmt.Errorf("at least one author ID is required")
	}

	// Verificar existência dos autores
	for _, authorID := range req.AuthorIDs {
		var exists bool
		err := r.DB.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM authors WHERE id = $1)`, authorID).Scan(&exists)
		if err != nil {
			log.Printf("Failed to check author existence: %v", err)
			return fmt.Errorf("failed to verify author")
		}
		if !exists {
			return fmt.Errorf("author with ID %d not found", authorID)
		}
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Inserir livro
	var book domain.Book
	err = tx.QueryRow(ctx, `
        INSERT INTO books (title, description) 
        VALUES ($1, $2) 
        RETURNING id, uuid, created_at`,
		req.Title, req.Description,
	).Scan(&book.ID, &book.UUID, &book.CreatedAt)

	if err != nil {
		log.Printf("Failed to insert book: %v", err)
		return fmt.Errorf("failed to insert book")
	}

	// Processar autores (todos já verificados)
	for _, authorID := range req.AuthorIDs {
		// Criar relação livro-autor
		_, err = tx.Exec(ctx, `
            INSERT INTO books_authors (book_id, author_id)
            VALUES ($1, $2)`,
			book.ID, authorID,
		)
		if err != nil {
			log.Printf("Failed to create books_author relation: %v", err)
			return fmt.Errorf("failed to create books-author relation")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return fmt.Errorf("failed to save data")
	}

	return nil
}

// GetBooksBasic retrieves books without author information (optimized)
func (r *BookRepository) GetBooksBasic(ctx context.Context, filters domain.BookFilters) ([]domain.Book, int, error) {
	// Build query
	query := `
        SELECT id, uuid, title, description, created_at
        FROM books
        WHERE ($1 = '' OR title ILIKE '%' || $1 || '%')
        ORDER BY ` + filters.Sort + ` ` + filters.SortDirection + `
        LIMIT $2 OFFSET $3`

	// Execute query
	rows, err := r.DB.Query(ctx, query, filters.Title, filters.Limit, filters.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// Process results
	var books []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.UUID, &b.Title, &b.Description, &b.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}
		books = append(books, b)
	}

	// Get total count (optimized count query)
	var total int
	countQuery := `SELECT COUNT(*) FROM books WHERE ($1 = '' OR title ILIKE '%' || $1 || '%')`
	if err := r.DB.QueryRow(ctx, countQuery, filters.Title).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	return books, total, nil
} // GetBooksWithAuthors retrieves books with author information (optimized join)
func (r *BookRepository) GetBooksWithAuthors(ctx context.Context, filters domain.BookFilters) ([]domain.Book, int, error) {
	// Build query
	query := `
        WITH paginated_books AS (
            SELECT id FROM books
            WHERE ($1 = '' OR title ILIKE '%' || $1 || '%')
            ORDER BY ` + filters.Sort + ` ` + filters.SortDirection + `
            LIMIT $2 OFFSET $3
        )
        SELECT 
            b.id, b.uuid, b.title, b.description, b.created_at,
            a.id as "author_id", a.uuid as "author_uuid", 
            a.name as "author_name", a.created_at as "author_created_at"
        FROM paginated_books pb
        JOIN books b ON pb.id = b.id
        LEFT JOIN books_authors ba ON b.id = ba.book_id
        LEFT JOIN authors a ON a.id = ba.author_id
        WHERE ($4 = '' OR a.name ILIKE '%' || $4 || '%')
        ORDER BY b.` + filters.Sort + ` ` + filters.SortDirection + `, a.name`

	// Execute query
	rows, err := r.DB.Query(ctx, query,
		filters.Title,
		filters.Limit,
		filters.Offset,
		filters.AuthorName,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// Process results with authors
	booksMap := make(map[int]*domain.Book)
	for rows.Next() {
		var b domain.Book
		var a domain.Author

		err := rows.Scan(
			&b.ID, &b.UUID, &b.Title, &b.Description, &b.CreatedAt,
			&a.ID, &a.UUID, &a.Name, &a.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}

		if _, exists := booksMap[b.ID]; !exists {
			booksMap[b.ID] = &b
			booksMap[b.ID].Authors = []*domain.Author{}
		}

		if a.ID != 0 { // Only add if author exists
			booksMap[b.ID].Authors = append(booksMap[b.ID].Authors, &a)
		}
	}

	// Convert map to slice
	books := make([]domain.Book, 0, len(booksMap))
	for _, book := range booksMap {
		books = append(books, *book)
	}

	// Get total count (with same filters)
	countQuery := `
        SELECT COUNT(DISTINCT b.id)
        FROM books b
        LEFT JOIN books_authors ba ON b.id = ba.book_id
        LEFT JOIN authors a ON a.id = ba.author_id
        WHERE ($1 = '' OR b.title ILIKE '%' || $1 || '%')
        AND ($2 = '' OR a.name ILIKE '%' || $2 || '%')`

	var total int
	if err := r.DB.QueryRow(ctx, countQuery, filters.Title, filters.AuthorName).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count failed: %w", err)
	}

	return books, total, nil
}
