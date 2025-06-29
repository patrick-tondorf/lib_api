package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/patrick-tondorf/lib_api/internal/domain"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	// Query atualizada para usar os campos corretos
	_, err := repo.db.Exec(
		ctx,
		`INSERT INTO users (email, password_hash) 
         VALUES ($1, $2)`,
		user.Email,
		user.PasswordHash, // Usando o hash, não a senha em texto puro
	)

	if err != nil {
		// Tratamento mais específico de erros
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return fmt.Errorf("user with this email already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := repo.db.QueryRow(
		ctx,
		`SELECT id, email, password_hash, created_at, updated_at 
         FROM users 
         WHERE email = $1`,
		email,
	)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash, // Agora pegando o hash, não a senha em texto puro
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}
