package domain

import "time"

type User struct {
	ID           int       `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Name         string    `db:"name" json:"name"`
	Role         string    `db:"role" json:"role"` // viewer, organizer, admin
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// Для регистрации и входа
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"` // viewer или organizer
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
