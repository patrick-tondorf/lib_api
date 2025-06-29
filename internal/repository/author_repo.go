package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/patrick-tondorf/lib_api/internal/domain"

	"github.com/jackc/pgx/v5"
)

type AuthorRepository struct {
	DB *pgx.Conn
}

func NewAuthorRepository(db *pgx.Conn) *AuthorRepository {
	return &AuthorRepository{DB: db}
}

// Crete
func (r *AuthorRepository) CreateAuthor(ctx context.Context, author *domain.Author) error {
	query := `
        INSERT INTO authors (name) 
        VALUES ($1)
        RETURNING id, uuid, created_at`

	err := r.DB.QueryRow(ctx, query, author.Name).
		Scan(&author.ID, &author.UUID, &author.CreatedAt)

	if err != nil {
		log.Printf("Error creating author: %v\n", err)
		return fmt.Errorf("failed to create author: %w", err)
	}

	log.Printf("Successfully created author with ID: %d\n", author.ID)
	return nil
}

// Get All
func (r *AuthorRepository) GetAuthors(ctx context.Context) ([]domain.Author, error) {
	log.Println("Attempting to query authors from database")

	query := `
        SELECT id, uuid, name, created_at, updated_at 
        FROM authors
        ORDER BY name`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		log.Printf("Database query error: %v\n", err)
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer rows.Close()

	var authors []domain.Author

	for rows.Next() {
		var a domain.Author
		err := rows.Scan(&a.ID, &a.UUID, &a.Name, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			log.Printf("Row scan error: %v\n", err)
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		authors = append(authors, a)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v\n", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	log.Printf("Successfully retrieved %d authors\n", len(authors))
	return authors, nil
}

func (r *AuthorRepository) GetAuthorsWithBooks(ctx context.Context) ([]domain.Author, error) {
	log.Println("Attempting to query authors from database")

	// Primeiro: buscar todos os autores
	authorsQuery := `
        SELECT id, uuid, name, created_at, updated_at 
        FROM authors
        ORDER BY name`

	authorRows, err := r.DB.Query(ctx, authorsQuery)
	if err != nil {
		log.Printf("Database query error for authors: %v\n", err)
		return nil, fmt.Errorf("database query error: %w", err)
	}
	defer authorRows.Close()

	var authors []domain.Author
	authorIDs := []int{} // Para coletar IDs dos autores encontrados

	for authorRows.Next() {
		var a domain.Author
		err := authorRows.Scan(&a.ID, &a.UUID, &a.Name, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			log.Printf("Row scan error for authors: %v\n", err)
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		authors = append(authors, a)
		authorIDs = append(authorIDs, a.ID)
	}

	if err := authorRows.Err(); err != nil {
		log.Printf("Rows error for authors: %v\n", err)
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// Se não encontrou autores, retorna vazio
	if len(authors) == 0 {
		return []domain.Author{}, nil
	}

	// Segundo: buscar todos os livros para os autores encontrados
	booksQuery := `
        SELECT b.id, b.uuid, b.title, b.description, b.created_at, ba.author_id
        FROM books b
        JOIN book_authors ba ON b.id = ba.book_id
        WHERE ba.author_id = ANY($1)
        ORDER BY ba.author_id, b.title`

	bookRows, err := r.DB.Query(ctx, booksQuery, authorIDs)
	if err != nil {
		log.Printf("Database query error for books: %v\n", err)
		return nil, fmt.Errorf("database books query error: %w", err)
	}
	defer bookRows.Close()

	// Criar um mapa de autores por ID para fácil associação
	authorMap := make(map[int]*domain.Author)
	for i := range authors {
		authorMap[authors[i].ID] = &authors[i]
		authors[i].Books = []*domain.Book{} // Inicializa slice vazia
	}

	// Processar livros e associar aos autores
	for bookRows.Next() {
		var b domain.Book
		var authorID int
		err := bookRows.Scan(&b.ID, &b.UUID, &b.Title, &b.Description, &b.CreatedAt, &authorID)
		if err != nil {
			log.Printf("Row scan error for books: %v\n", err)
			return nil, fmt.Errorf("book row scan error: %w", err)
		}

		if author, exists := authorMap[authorID]; exists {
			author.Books = append(author.Books, &b)
		}
	}

	if err := bookRows.Err(); err != nil {
		log.Printf("Rows error for books: %v\n", err)
		return nil, fmt.Errorf("book rows error: %w", err)
	}

	log.Printf("Successfully retrieved %d authors with their books\n", len(authors))
	return authors, nil
}

func (r *AuthorRepository) GetAuthorByID(ctx context.Context, id int) (*domain.Author, error) {
	log.Printf("Attempting to query author with ID: %d\n", id)

	// Busca o autor
	author := &domain.Author{}
	err := r.DB.QueryRow(ctx, `
        SELECT id, uuid, name, created_at, updated_at 
        FROM authors 
        WHERE id = $1`, id).
		Scan(&author.ID, &author.UUID, &author.Name, &author.CreatedAt, &author.UpdatedAt)

	if err != nil {
		log.Printf("Error fetching author: %v\n", err)
		return nil, fmt.Errorf("error fetching author: %w", err)
	}

	// Busca os livros do autor
	rows, err := r.DB.Query(ctx, `
        SELECT b.id, b.uuid, b.title, b.description, b.created_at
        FROM books b
        JOIN book_authors ba ON b.id = ba.book_id
        WHERE ba.author_id = $1
        ORDER BY b.title`, id)

	if err != nil {
		log.Printf("Error fetching author's books: %v\n", err)
		return nil, fmt.Errorf("error fetching author's books: %w", err)
	}
	defer rows.Close()

	author.Books = []*domain.Book{}
	for rows.Next() {
		var b domain.Book
		err := rows.Scan(&b.ID, &b.UUID, &b.Title, &b.Description, &b.CreatedAt)
		if err != nil {
			log.Printf("Error scanning book: %v\n", err)
			continue // Ou retornar erro se preferir
		}
		author.Books = append(author.Books, &b)
	}

	log.Printf("Successfully retrieved author with %d books\n", len(author.Books))
	return author, nil
}
