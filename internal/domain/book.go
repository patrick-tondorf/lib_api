package domain

import "time"

type Book struct {
	ID        int       `json:"id" example:"1"`                 // @example 1
	Title     string    `json:"title" example:"1984"`           // @example 1984
	Author    string    `json:"author" example:"George Orwell"` // @example George Orwell
	CreatedAt time.Time `json:"created_at,omitempty"`
} //@name Book
