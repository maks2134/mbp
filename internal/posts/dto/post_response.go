package dto

import "time"

type PostResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Tag          string    `json:"tag"`
	Like         int       `json:"like"`
	CountViewers int       `json:"count_viewers"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
