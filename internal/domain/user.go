package domain

import "time"

type User struct {
	ID           string     `json:"-" db:"id"`
	UUID         string     `json:"-" db:"uuid"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"password" db:"-"`                                  // Usado apenas para receber o input
	PasswordHash string     `json:"-" db:"password_hash" swaggerignore:"true"`        //swagger:ignore
	CreatedAt    time.Time  `json:"-" db:"created_at"  swaggerignore:"true`           //swagger:ignore
	UpdatedAt    *time.Time `json:"-,omitempty" db:"updated_at"  swaggerignore:"true` //swagger:ignore
}

type UserListResponse struct {
	ID           string     `json:"id" db:"id"`
	UUID         string     `json:"uuid" db:"uuid"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"` // Armazenado no banco
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
}
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
} //@name Credential
