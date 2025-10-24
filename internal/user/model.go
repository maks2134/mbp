package model

import "time"

type User struct {
	ID           int        `db:"id"`
	Name         string     `db:"name"`
	Username     string     `db:"username"`
	PasswordHash string     `db:"password_hash"`
	Email        *string    `db:"email"`
	Age          int        `db:"age"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}
