package user

import "time"

type User struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name" validator:"required"`
	Age       int        `db:"age" json:"age" validator:"required"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
