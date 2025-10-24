package auth

import "time"

type Auth struct {
	ID             int        `db:"id" json:"id"`
	UserID         int        `db:"user_id" json:"user_id"`
	Email          string     `db:"email" json:"email"`
	PasswordHash   string     `db:"password_hash" json:"password"`
	LastLoginAt    *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
	FailedAttempts int        `db:"failed_attempts" json:"failed_attempts"`
	LockedUntil    *time.Time `db:"locked_until" json:"locked_until,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}
