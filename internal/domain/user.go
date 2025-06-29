package domain

import "time"

type User struct {
	ID           string     `json:"-" db:"id"`
	UUID         string     `json:"-" db:"uuid"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"password" db:"-"`      // Usado apenas para receber o input
	PasswordHash string     `json:"-" db:"password_hash"` // Armazenado no banco
	CreatedAt    time.Time  `json:"-" db:"created_at"`
	UpdatedAt    *time.Time `json:"-,omitempty" db:"updated_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
} //@name Credential
