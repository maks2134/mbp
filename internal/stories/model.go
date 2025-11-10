package stories

import "time"

type Story struct {
	ID         int        `db:"id" json:"id"`
	UserID     int        `db:"user_id" json:"user_id"`
	FileURL    string     `db:"file_url" json:"file_url"`
	FileType   string     `db:"file_type" json:"file_type"`
	ViewsCount int        `db:"views_count" json:"views_count"`
	ExpiresAt  time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type StoryView struct {
	ID       int       `db:"id" json:"id"`
	StoryID  int       `db:"story_id" json:"story_id"`
	UserID   int       `db:"user_id" json:"user_id"`
	ViewedAt time.Time `db:"viewed_at" json:"viewed_at"`
}
