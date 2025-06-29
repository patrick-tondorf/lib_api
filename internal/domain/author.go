package domain

import "time"

type Author struct {
	ID        int        `json:"-" swaggerignore:"true" db:"id"`
	UUID      string     `json:"uuid" swaggerignore:"true" db:"uuid"` // Ignora no input
	Name      string     `json:"name" example:"George Orwell" db:"name"`
	Bio       string     `json:"bio,omitempty" example:"Autor de 1984 e A Revolução dos Bichos" db:"bio"` // Biografia do autor
	Books     []*Book    `json:"books,omitempty" swaggerignore:"true"`                                    // Lista de livros do autor
	CreatedAt time.Time  `json:"createdAt" swaggerignore:"true" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"  swaggerignore:"true" db:"updated_at"`
} // @name Author
