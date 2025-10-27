package posts

import (
	"time"
)

type Tag string

type Post struct {
	ID           int        `db:"id"`
	UserID       int        `db:"user_id"`
	Title        string     `db:"title"`
	Description  string     `db:"description"`
	Tag          string     `db:"tag"`
	Like         int        `db:"like"`
	CountViewers int        `db:"count_viewers"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}
